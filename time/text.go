package main

import (
	"fmt"
	"time"
)

var (
	hour   = 4
	minute = 30
)

func main() {
	nextRun := calculateNextRun(time.Now())
	//timer := time.NewTimer(time.Until(nextRun))
	fmt.Printf("Next run scheduled at: %s\n", nextRun)
}

func calculateNextRun(now time.Time) time.Time {
	// Today's scheduled time
	scheduled := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		hour,
		minute,
		0, 0,
		now.Location(),
	)

	// If today's scheduled time hasn't passed yet
	if now.Before(scheduled) {
		return scheduled
	}

	// Otherwise schedule for tomorrow
	return scheduled.Add(24 * time.Hour)
}
