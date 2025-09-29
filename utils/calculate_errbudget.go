package utils

import (
	"fmt"
	"slo-tracker/schema"
	"strconv"
	"strings"
	"time"
)

func CalculateErrBudget(targetSLOinPerc float32) float32 {
	totalSecsInYear := 31536000
	downtimeInFraction := 1 - (targetSLOinPerc / 100)
	errBudgetInMin := (downtimeInFraction * float32(totalSecsInYear)) / 60
	return errBudgetInMin
}

func CalculateMonthlyErrBudget(SLO *schema.SLO, incidents []*schema.Incident, yearMonth string, schedule []schema.StoreWorkingSchedule) (float32, float32) {
	year, month, _ := ParseYearMonth(yearMonth)

	// First day of the given month
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// First day of the next month
	end := start.AddDate(0, 1, 0)

	// Duration in minutes
	minutesInMonth := end.Sub(start).Minutes() * float64(SLO.TargetSLO) / 100

	var minutesInAlarm float32 = 0.0

	for _, incident := range incidents {
		if incident.MarkFalsePositive {
			continue
		}

		m, _ := DowntimeAcrossDays(*incident.CreatedAt, incident.RealErrorBudget, schedule)
		minutesInAlarm += m
	}

	remainingErrBudget := float32(minutesInMonth) - minutesInAlarm

	return remainingErrBudget, 100 - (minutesInAlarm * 100 / float32(minutesInMonth))
}

func ParseYearMonth(yearMonth string) (int, int, error) {
	parts := strings.Split(yearMonth, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format, expected YYYY-MM")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year: %w", err)
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid month: %w", err)
	}

	if month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("month must be between 1 and 12, got %d", month)
	}

	return year, month, nil
}

func DowntimeAcrossDays(alarmStart time.Time, durationMinutes float32, schedule []schema.StoreWorkingSchedule) (float32, error) {

	alarmEnd := alarmStart.Add(time.Duration(float64(durationMinutes) * float64(time.Minute)))
	var totalMinutes float64

	currentDay := alarmStart

	for !currentDay.After(alarmEnd) {

		weekday := int(currentDay.Weekday())
		var daySchedule *schema.StoreWorkingSchedule

		// Find the schedule for this weekday
		for _, s := range schedule {
			if s.Weekday == weekday {
				daySchedule = &s
				break
			}
		}

		if daySchedule != nil {
			// Parse open/close hours
			openParts := strings.Split(daySchedule.OpenHour, ":")
			closeParts := strings.Split(daySchedule.CloseHour, ":")

			openHour, _ := time.ParseDuration(fmt.Sprintf("%sh%sm%ss", openParts[0], openParts[1], openParts[2]))
			closeHour, _ := time.ParseDuration(fmt.Sprintf("%sh%sm%ss", closeParts[0], closeParts[1], closeParts[2]))

			// Create actual time.Time for open/close today
			dayOpen := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), 0, 0, 0, 0, currentDay.Location()).Add(openHour)
			dayClose := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), 0, 0, 0, 0, currentDay.Location()).Add(closeHour)

			// Calculate overlap
			start := maxTime(currentDay, dayOpen)
			end := minTime(alarmEnd, dayClose)

			if end.After(start) {
				totalMinutes += end.Sub(start).Minutes()
			}
		}

		currentDay = currentDay.AddDate(0, 0, 1)
	}

	return float32(totalMinutes), nil
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
