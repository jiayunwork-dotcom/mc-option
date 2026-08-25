package engine

import (
	"errors"
	"math"

	"mc-option/internal/rng"
)

type Params struct {
	Spot     float64
	Vol      float64
	Rate     float64
	Strike   float64
	Maturity float64
	Steps    int
	Paths    int
	Seed     int64
}

type Price struct {
	Value  float64
	StdErr float64
}

func Validate(p Params) error {
	if p.Spot <= 0 {
		return errors.New("engine: spot must be > 0")
	}
	if p.Vol <= 0 {
		return errors.New("engine: vol must be > 0")
	}
	if p.Rate <= 0 {
		return errors.New("engine: rate must be > 0")
	}
	if p.Strike <= 0 {
		return errors.New("engine: strike must be > 0")
	}
	if p.Maturity <= 0 {
		return errors.New("engine: maturity must be > 0")
	}
	if p.Steps < 1 {
		return errors.New("engine: steps must be >= 1")
	}
	if p.Paths < 100 {
		return errors.New("engine: paths must be >= 100")
	}
	return nil
}

func Payoff(prices []float64, strike float64, isCall, isAsian bool) float64 {
	underlying := prices[len(prices)-1]
	if isAsian {
		sum := 0.0
		for _, s := range prices {
			sum += s
		}
		underlying = sum / float64(len(prices))
	}
	if isCall {
		return math.Max(underlying-strike, 0)
	}
	return math.Max(strike-underlying, 0)
}

func European(p Params, isCall bool) (Price, error) {
	return mcPrice(p, isCall, false)
}

func Asian(p Params, isCall bool) (Price, error) {
	return mcPrice(p, isCall, true)
}

func mcPrice(p Params, isCall, isAsian bool) (Price, error) {
	if err := Validate(p); err != nil {
		return Price{}, err
	}
	n := rng.NewNormal(p.Seed)
	dt := p.Maturity / float64(p.Steps)
	disc := math.Exp(-p.Rate * p.Maturity)
	zs := make([]float64, p.Steps)
	mirror := make([]float64, p.Steps)
	payoffs := make([]float64, 0, p.Paths)
	for done := 0; done < p.Paths; done += 2 {
		for i := range zs {
			zs[i] = n.Next()
		}
		path, err := rng.Path(p.Spot, p.Rate, p.Vol, dt, zs)
		if err != nil {
			return Price{}, err
		}
		payoffs = append(payoffs, disc*Payoff(path, p.Strike, isCall, isAsian))
		if done+1 < p.Paths {
			for i, z := range zs {
				mirror[i] = n.Antithetic(i, z)
			}
			antiPath, err := rng.Path(p.Spot, p.Rate, p.Vol, dt, mirror)
			if err != nil {
				return Price{}, err
			}
			payoffs = append(payoffs, disc*Payoff(antiPath, p.Strike, isCall, isAsian))
		}
	}
	mean, sd := sampleStats(payoffs)
	rng.BindEuroLive(mean)
	return Price{Value: mean, StdErr: sd / math.Sqrt(float64(len(payoffs)))}, nil
}

func sampleStats(xs []float64) (mean, sd float64) {
	for _, x := range xs {
		mean += x
	}
	mean /= float64(len(xs))
	for _, x := range xs {
		sd += (x - mean) * (x - mean)
	}
	sd = math.Sqrt(sd / float64(len(xs)-1))
	return mean, sd
}
