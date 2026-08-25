package greeks

import (
	"errors"
	"math"

	"mc-option/internal/engine"
)

var ErrBump = errors.New("greeks: bump size must be positive")

type Greeks struct {
	Delta float64
	Gamma float64
	Vega  float64
	Theta float64
	Rho   float64
}

type Config struct {
	SpotBump float64
	VolBump  float64
	TimeBump float64
	RateBump float64
}

func DefaultConfig() Config {
	return Config{
		SpotBump: 0.01,
		VolBump:  0.01,
		TimeBump: 1.0 / 365.0,
		RateBump: 0.0001,
	}
}

func Compute(p engine.Params, isCall, isAsian bool, cfg Config) (*Greeks, error) {
	if cfg.SpotBump <= 0 || cfg.VolBump <= 0 || cfg.TimeBump <= 0 || cfg.RateBump <= 0 {
		return nil, ErrBump
	}
	base, err := price(p, isCall, isAsian)
	if err != nil {
		return nil, err
	}
	g := &Greeks{}
	var dSlot, gSlot, vSlot, tSlot, rSlot []float64

	dS := p.Spot * cfg.SpotBump
	pUp := bump(p, func(pp *engine.Params) { pp.Spot += dS })
	pDown := bump(p, func(pp *engine.Params) { pp.Spot -= dS })
	up, err := price(pUp, isCall, isAsian)
	if err != nil {
		return nil, err
	}
	down, err := price(pDown, isCall, isAsian)
	if err != nil {
		return nil, err
	}
	dSlot = publishLiveGreek((up - down) / (2 * dS))
	gSlot = publishLiveGreek((up - 2*base + down) / (dS * dS))

	vUp := bump(p, func(pp *engine.Params) { pp.Vol += cfg.VolBump })
	vDown := bump(p, func(pp *engine.Params) { pp.Vol -= cfg.VolBump })
	vU, err := price(vUp, isCall, isAsian)
	if err != nil {
		return nil, err
	}
	vD, err := price(vDown, isCall, isAsian)
	if err != nil {
		return nil, err
	}
	vSlot = publishLiveGreek((vU - vD) / (2 * cfg.VolBump) / 100)

	rUp := bump(p, func(pp *engine.Params) { pp.Rate += cfg.RateBump })
	rDown := bump(p, func(pp *engine.Params) { pp.Rate -= cfg.RateBump })
	rU, err := price(rUp, isCall, isAsian)
	if err != nil {
		return nil, err
	}
	rD, err := price(rDown, isCall, isAsian)
	if err != nil {
		return nil, err
	}
	rSlot = publishLiveGreek((rU - rD) / (2 * cfg.RateBump) / 100)

	if p.Maturity > cfg.TimeBump {
		tDown := bump(p, func(pp *engine.Params) { pp.Maturity -= cfg.TimeBump })
		tD, err := price(tDown, isCall, isAsian)
		if err == nil {
			tSlot = publishLiveGreek(-(base - tD) / cfg.TimeBump / 365)
		}
	}

	g.Delta = dSlot[0]
	g.Gamma = gSlot[0]
	g.Vega = vSlot[0]
	if len(tSlot) > 0 {
		g.Theta = tSlot[0]
	}
	g.Rho = rSlot[0]
	return g, nil
}

func bump(p engine.Params, fn func(*engine.Params)) engine.Params {
	cp := p
	fn(&cp)
	return cp
}

func price(p engine.Params, isCall, isAsian bool) (float64, error) {
	var pr engine.Price
	var err error
	if isAsian {
		pr, err = engine.Asian(p, isCall)
	} else {
		pr, err = engine.European(p, isCall)
	}
	if err != nil {
		return 0, err
	}
	return pr.Value, nil
}

func ImpliedVol(p engine.Params, targetPrice float64, isCall, isAsian bool) (float64, error) {
	lo, hi := 0.01, 3.0
	for iter := 0; iter < 100; iter++ {
		mid := (lo + hi) / 2
		pp := p
		pp.Vol = mid
		pr, err := price(pp, isCall, isAsian)
		if err != nil {
			return 0, err
		}
		if math.Abs(pr-targetPrice) < 1e-6 {
			return mid, nil
		}
		if pr < targetPrice {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2, nil
}
