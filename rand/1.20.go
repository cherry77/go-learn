package main

import (
	"math"
	"math/rand"
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
func GenerateNormallyDistributedTimeV_1_20(days int, sigma, meanDay float64) time.Time {
	now := time.Now()
	startTime := now.AddDate(0, 0, -days)
	timeWindow := float64(now.Sub(startTime))

	// Normalize mean to [0,1] range
	normalizedMean := meanDay / float64(days)

	// Generate and adjust normally distributed value
	// Go 1.20 uses math/rand instead of math/rand/v2
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	adjustedValue := (rng.NormFloat64() * sigma) + normalizedMean

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
