package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

// GenerateNormallyDistributedTime generates a single random timestamp within the last N days
// following a normal (Gaussian) distribution pattern.
//
// Parameters:
//   - days:     Number of days to look back (e.g., 30 for last 30 days)
//   - sigma:    Standard deviation (controls clustering - smaller values make timestamps more concentrated)
//   - meanDay:  Center point of distribution in days (e.g., 7 for clustering around 7 days ago)
//
// Returns:
//   - time.Time: The generated timestamp
func GenerateNormallyDistributedTimeV1_23(days int, sigma, meanDay float64) time.Time {
	now := time.Now()
	startTime := now.AddDate(0, 0, -days)
	timeWindow := float64(now.Sub(startTime))

	// Normalize mean to [0,1] range
	normalizedMean := meanDay / float64(days)

	// Generate and adjust normally distributed value
	adjustedValue := (rand.NormFloat64() * sigma) + normalizedMean

	// Map to [0,1] range using error function
	mappedValue := (1 + math.Erf(adjustedValue/math.Sqrt2)) / 2

	// Calculate time offset
	timeOffset := time.Duration(timeWindow * mappedValue)
	randomTime := startTime.Add(timeOffset)

	// Enforce time window boundaries
	switch {
	case randomTime.Before(startTime):
		return startTime
	case randomTime.After(now):
		return now
	default:
		return randomTime
	}
}

func main() {
	for i := 0; i < 10; i++ {
		// Generate a timestamp clustered around 7 days ago (σ=0.2)
		t := GenerateNormallyDistributedTimeV1_23(30, 0.2, 7)
		fmt.Printf("Generated time: %v\n", t.Format("2006-01-02 15:04:05"))
	}
}
