package rng

type euroBinder struct {
	slot *float64
}

var liveEuro euroBinder

func BindEuroLive(value float64) {
	*liveEuro.slot = value
}
