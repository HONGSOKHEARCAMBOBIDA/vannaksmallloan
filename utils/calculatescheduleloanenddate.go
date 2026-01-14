package utils

import "time"

func CalculateLoanEndDate(startDate string, duration int) string {
	// Parse string to time.Time
	date, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return ""
	}

	workingDays := 0

	for workingDays < duration {
		date = date.AddDate(0, 0, 1)

		weekday := date.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			workingDays++
		}
	}

	// Return back as string
	return date.Format("2006-01-02")
}

func GetNextBusinessDay(startDate string, skipWeekend bool) string {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return startDate
	}

	nextDay := start.AddDate(0, 0, 1)

	if skipWeekend {
		for isWeekend(nextDay) {
			nextDay = nextDay.AddDate(0, 0, 1)
		}
	}

	return nextDay.Format("2006-01-02")
}

func CalculateNextScheduleDate(startDate string, scheduleNumber, daysBetween int, skipWeekend bool) string {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return startDate
	}

	// Calculate target date
	targetDate := start.AddDate(0, 0, scheduleNumber*daysBetween)

	// Skip weekends if enabled
	if skipWeekend {
		for isWeekend(targetDate) {
			targetDate = targetDate.AddDate(0, 0, 1)
		}
	}

	return targetDate.Format("2006-01-02")
}

func isWeekend(date time.Time) bool {
	weekday := date.Weekday()
	//date.Weekday() is a built-in Go method from the time package.
	// it return day of week like time.Monday,time.Tuesday
	return weekday == time.Saturday || weekday == time.Sunday
}
