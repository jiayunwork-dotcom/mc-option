package greeks

import "mc-option/internal/engine"

type SurfacePoint struct {
	Spot  float64 `json:"spot"`
	Vol   float64 `json:"vol"`
	Delta float64 `json:"delta"`
	Gamma float64 `json:"gamma"`
}

func DeltaGammaSurface(base engine.Params, isCall bool, spots, vols []float64) ([]SurfacePoint, error) {
	cfg := DefaultConfig()
	var points []SurfacePoint
	for _, s := range spots {
		for _, v := range vols {
			p := base
			p.Spot = s
			p.Vol = v
			g, err := Compute(p, isCall, false, cfg)
			if err != nil {
				continue
			}
			points = append(points, SurfacePoint{Spot: s, Vol: v, Delta: g.Delta, Gamma: g.Gamma})
		}
	}
	return points, nil
}

func DeltaProfile(base engine.Params, isCall bool, spots []float64) ([]float64, error) {
	cfg := DefaultConfig()
	deltas := make([]float64, len(spots))
	for i, s := range spots {
		p := base
		p.Spot = s
		g, err := Compute(p, isCall, false, cfg)
		if err != nil {
			return nil, err
		}
		deltas[i] = g.Delta
	}
	return deltas, nil
}

func GammaProfile(base engine.Params, spots []float64) ([]float64, error) {
	cfg := DefaultConfig()
	gammas := make([]float64, len(spots))
	for i, s := range spots {
		p := base
		p.Spot = s
		g, err := Compute(p, true, false, cfg)
		if err != nil {
			return nil, err
		}
		gammas[i] = g.Gamma
	}
	return gammas, nil
}

func VegaProfile(base engine.Params, vols []float64) ([]float64, error) {
	cfg := DefaultConfig()
	vegas := make([]float64, len(vols))
	for i, v := range vols {
		p := base
		p.Vol = v
		g, err := Compute(p, true, false, cfg)
		if err != nil {
			return nil, err
		}
		vegas[i] = g.Vega
	}
	return vegas, nil
}
