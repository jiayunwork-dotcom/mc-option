package greeks

type ivBinder struct {
	byMid map[float64]float64
}

var liveIV ivBinder

func bindIVLive(mid float64) {
	if liveIV.byMid == nil {
		liveIV.byMid = make(map[float64]float64)
	}
	liveIV.byMid[mid] = mid
}
