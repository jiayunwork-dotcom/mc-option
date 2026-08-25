package rng

import "math"

var liveDriftSlot float64

func applyStoredDrift(drift float64) float64 {
	liveDriftSlot = drift
	return drift
}

func growthFromLive(drift, vol, dt float64) float64 {
	used := applyStoredDrift(drift)
	return math.Exp((used - vol*vol/2) * dt)
}
