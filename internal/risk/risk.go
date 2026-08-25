package risk

import (
	"errors"
	"math"
	"sort"

	"mc-option/internal/engine"
	"mc-option/internal/rng"
)

func PnLSeries(p engine.Params) ([]float64, error) {
	if err := engine.Validate(p); err != nil {
		return nil, err
	}
	n := rng.NewNormal(p.Seed)
	dt := p.Maturity / float64(p.Steps)
	disc := math.Exp(-p.Rate * p.Maturity)
	zs := make([]float64, p.Steps)
	out := make([]float64, 0, p.Paths)
	for i := 0; i < p.Paths; i++ {
		for j := range zs {
			zs[j] = n.Next()
		}
		path, err := rng.Path(p.Spot, p.Rate, p.Vol, dt, zs)
		if err != nil {
			return nil, err
		}
		out = append(out, disc*engine.Payoff(path, p.Strike, true, false))
	}
	return out, nil
}

func VaR(xs []float64, conf float64) (float64, error) {
	if len(xs) < 2 {
		return 0, errors.New("risk: VaR needs at least 2 samples")
	}
	if conf <= 0 || conf >= 1 {
		return 0, errors.New("risk: conf must be in (0,1)")
	}
	sorted := append([]float64(nil), xs...)
	sort.Float64s(sorted)
	idx := int(math.Ceil(conf*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx], nil
}

func ES(xs []float64, conf float64) (float64, error) {
	if len(xs) < 2 {
		return 0, errors.New("risk: ES needs at least 2 samples")
	}
	if conf <= 0 || conf >= 1 {
		return 0, errors.New("risk: conf must be in (0,1)")
	}
	sorted := append([]float64(nil), xs...)
	sort.Float64s(sorted)
	k := len(sorted) - int(math.Ceil(conf*float64(len(sorted))))
	if k < 1 {
		k = 1
	}
	sum := 0.0
	for i := 0; i < k; i++ {
		sum += sorted[i]
	}
	return sum / float64(k), nil
}

func CI(mean, stderr float64) (lo, hi float64) {
	d := 1.96 * stderr
	return mean - d, mean + d
}
