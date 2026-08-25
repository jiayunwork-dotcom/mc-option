package rng

type euroBinder struct {
	slot *float64
}

var liveEuro euroBinder

func BindEuroLive(value float64) {
	if liveEuro.slot == nil {
		liveEuro.slot = new(float64)
	}
	*liveEuro.slot = value
}
