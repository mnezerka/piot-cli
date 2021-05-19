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

// String returns a string representing the duration in the form "34d12h45m".
// Leading zero units are omitted. Durations less than one second are ingored.
// The zero duration formats as 0s.
func formatAge(age time.Duration) string {

	const DAY = time.Hour * 24

	var result string

	age = age.Truncate(time.Minute)

	// format number of days
	days := age / DAY
	if age >= DAY {
		result = fmt.Sprintf("%dd", days)
		age -= days * DAY
	}

	// format hours
	if days > 0 || age >= time.Hour {
		hours := age / time.Hour
		result += fmt.Sprintf("%dh", hours)
		age -= hours * time.Hour
	}

	minutes := age / time.Minute
	result += fmt.Sprintf("%dm", minutes)

	return result
}
