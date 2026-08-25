package engine

type knockBinder struct {
	byValue map[float64]float64
}

var liveKnock knockBinder

func BindKnockLive(value float64) {
	if liveKnock.byValue == nil {
		liveKnock.byValue = make(map[float64]float64)
	}
	liveKnock.byValue[value] = value
}
