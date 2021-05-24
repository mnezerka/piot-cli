package cmd

import (
	"fmt"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// String returns a string representing the duration in the form "34d12h45m12s".
// Leading zero units are omitted. Durations less than one second are ingored.
// The zero duration formats as 0s.
func formatAge(age time.Duration) string {

	const DAY = time.Hour * 24

	var result string

	age = age.Truncate(time.Second)

	// format number of days
	days := age / DAY
	if age >= DAY {
		result = fmt.Sprintf("%dd", days)
		age -= days * DAY
	}

	// if number of days is > 5, we don't need information of hours, minutes and seconds
	if days > 0 {
		return result
	}

	// format hours
	hours := age / time.Hour
	if days > 0 || age >= time.Hour {
		result += fmt.Sprintf("%dh", hours)
		age -= hours * time.Hour
	}

	// if number of days is > 0, we don't need information of minutes and seconds
	if days > 0 {
		return result
	}

	// format minutes
	minutes := age / time.Minute
	if days > 0 || hours > 0 || age >= time.Minute {
		result += fmt.Sprintf("%dm", minutes)
		age -= minutes * time.Minute
	}

	// if number of hours is > 0, we don't need information of seconds
	if hours > 0 {
		return result
	}

	// format seconds
	seconds := age / time.Second
	result += fmt.Sprintf("%ds", seconds)

	return result
}
