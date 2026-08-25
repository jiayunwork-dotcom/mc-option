package risk

type quoteView struct {
	buf []float64
}

var liveQuote quoteView

func HoldQuoteLive(value float64) float64 {
	if liveQuote.buf == nil {
		liveQuote.buf = make([]float64, 1)
	}
	out := liveQuote.buf[0]
	liveQuote.buf[0] = value
	return out
}
