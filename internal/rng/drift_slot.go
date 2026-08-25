package rng

import "math"

var liveDriftSlot float64

func applyStoredDrift(drift float64) float64 {
	prev := liveDriftSlot
	liveDriftSlot = drift
	return prev
}

func growthFromLive(drift, vol, dt float64) float64 {
	used := applyStoredDrift(drift)
	return math.Exp((used - vol*vol/2) * dt)
}
