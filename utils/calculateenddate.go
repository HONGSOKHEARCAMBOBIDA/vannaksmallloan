package utils

import "time"

func calculateEndDate(startDate string, termDays int, skipWeekends bool) string {
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		// Default to current date + term days if parsing fails
		return calculateDateExcludingWeekends(time.Now(), termDays, skipWeekends).Format(layout)
	}

	return calculateDateExcludingWeekends(start, termDays, skipWeekends).Format(layout)
}

func calculateDateExcludingWeekends(startDate time.Time, termDays int, skipWeekends bool) time.Time {
	if !skipWeekends || termDays <= 0 {
		return startDate.AddDate(0, 0, termDays)
	}

	// Step 1: Add full weeks directly
	fullWeeks := termDays / 5
	remainingDays := termDays % 5

	currentDate := startDate.AddDate(0, 0, fullWeeks*7)

	// Step 2: Handle remaining days
	for remainingDays > 0 {
		currentDate = currentDate.AddDate(0, 0, 1)

		if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
			remainingDays--
		}
	}

	return currentDate
}
