package utils

import (
	"time"
)

func CalculateErrBudget(targetSLOinPerc float32) float32 {
	totalSecsInYear := 31536000
	downtimeInFraction := 1 - (targetSLOinPerc / 100)
	errBudgetInMin := (downtimeInFraction * float32(totalSecsInYear)) / 60
	return errBudgetInMin
}

func NextWorkStart(t *time.Time) *time.Time {
	day := t.Truncate(24 * time.Hour)

	startOfDay := day.Add(time.Duration(9) * time.Hour)
	endOfDay := day.Add(time.Duration(18) * time.Hour)

	switch {
	case t.Before(startOfDay):
		// Before work start → same day at 9:00
		return &startOfDay
	case t.After(endOfDay):
		// After work end → next day at 9:00
		nextDayStart := day.Add(24 * time.Hour).Add(time.Duration(9) * time.Hour)
		return &nextDayStart
	default:
		// Within work hours → unchanged
		return t
	}
}

func CalculateAccountableErrBudget(alarmStart *time.Time, minutes float32) float32 {

	nextWorkStart := NextWorkStart(alarmStart)

	diffMinutes := nextWorkStart.Sub(*alarmStart)

	if float32(diffMinutes.Minutes()) > minutes {
		return 0
	}

	return minutes - float32(diffMinutes.Minutes())
}
