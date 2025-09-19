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

func CalculateMonthlyErrBudget(SLO *schema.SLO, incidents []*schema.Incident, yearMonth string) (float32, float32) {
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

		m, _ := DowntimeAcrossDays(SLO.OpenHour, SLO.CloseHour, *incident.CreatedAt, incident.ErrorBudgetSpent)
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

func DowntimeAcrossDays(openStr, closeStr string, alarmStart time.Time, durationMinutes float32) (float32, error) {
	// Parse open and close times (use today's date just for parsing HH:mm:ss)
	openParsed, err := time.Parse("15:04:05", openStr)
	if err != nil {
		return 0, fmt.Errorf("invalid open time: %w", err)
	}
	closeParsed, err := time.Parse("15:04:05", closeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid close time: %w", err)
	}

	alarmEnd := alarmStart.Add(time.Duration(float64(durationMinutes) * float64(time.Minute)))
	var totalMinutes float64

	// Start from the day of alarmStart
	currentDay := alarmStart

	for !currentDay.After(alarmEnd) {
		// Create open/close times anchored to currentDay
		open := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(),
			openParsed.Hour(), openParsed.Minute(), openParsed.Second(), 0, currentDay.Location())
		close := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(),
			closeParsed.Hour(), closeParsed.Minute(), closeParsed.Second(), 0, currentDay.Location())

		// Find overlap
		start := maxTime(open, alarmStart)
		end := minTime(close, alarmEnd)

		if end.After(start) {
			totalMinutes += end.Sub(start).Minutes()
		}

		// Next day
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
