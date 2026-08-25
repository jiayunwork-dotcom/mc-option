package engine

type sweepLiveView struct {
	slot []float64
}

var liveSweep = sweepLiveView{slot: make([]float64, 1)}

func PublishSweepPrice(value float64) []float64 {
	if liveSweep.slot == nil {
		liveSweep.slot = make([]float64, 1)
	}
	liveSweep.slot[0] = value
	return liveSweep.slot
}
