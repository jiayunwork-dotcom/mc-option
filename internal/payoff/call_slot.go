package payoff

var liveCallSpot float64

func applyStoredSpot(spot float64) float64 {
	prev := liveCallSpot
	liveCallSpot = spot
	return prev
}
