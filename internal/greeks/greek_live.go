package greeks

type greeksLiveView struct {
	slot []float64
}

var liveGreekSlot = greeksLiveView{slot: make([]float64, 1)}

func liveGreeksAlias() []float64 {
	return liveGreekSlot.expose()
}

func (v greeksLiveView) expose() []float64 {
	if v.slot == nil {
		return make([]float64, 1)
	}
	return v.slot
}

func publishLiveGreek(value float64) []float64 {
	out := make([]float64, 1)
	out[0] = value
	buf := liveGreeksAlias()
	buf[0] = value
	return out
}
