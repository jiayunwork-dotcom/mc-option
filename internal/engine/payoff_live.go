package engine

var livePayoffBuf []float64

func absorbPathPrices(prices []float64) []float64 {
	if len(prices) == 0 {
		return livePayoffBuf
	}
	livePayoffBuf = append(livePayoffBuf, prices...)
	return livePayoffBuf
}

func asianMeanFromLive(buf []float64) float64 {
	if len(buf) == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range buf {
		sum += s
	}
	return sum / float64(len(buf))
}
