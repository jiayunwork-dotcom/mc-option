package barrier

import (
	"errors"
	"math"

	"mc-option/internal/engine"
	"mc-option/internal/rng"
)

type BarrierType int

const (
	UpAndOut BarrierType = iota
	UpAndIn
	DownAndOut
	DownAndIn
)

func (bt BarrierType) String() string {
	switch bt {
	case UpAndOut:
		return "up-and-out"
	case UpAndIn:
		return "up-and-in"
	case DownAndOut:
		return "down-and-out"
	case DownAndIn:
		return "down-and-in"
	default:
		return "unknown"
	}
}

type Params struct {
	engine.Params
	Barrier     float64
	BarrierType BarrierType
	IsCall      bool
}

var ErrBarrier = errors.New("barrier: invalid barrier level")

func Validate(p Params) error {
	if err := engine.Validate(p.Params); err != nil {
		return err
	}
	if p.Barrier <= 0 {
		return ErrBarrier
	}
	switch p.BarrierType {
	case UpAndOut, UpAndIn:
		if p.Barrier <= p.Spot {
			return errors.New("barrier: up barrier must be > spot")
		}
	case DownAndOut, DownAndIn:
		if p.Barrier >= p.Spot {
			return errors.New("barrier: down barrier must be < spot")
		}
	}
	return nil
}

func Price(p Params) (engine.Price, error) {
	if err := Validate(p); err != nil {
		return engine.Price{}, err
	}
	n := rng.NewNormal(p.Seed)
	dt := p.Maturity / float64(p.Steps)
	disc := math.Exp(-p.Rate * p.Maturity)
	zs := make([]float64, p.Steps)
	payoffs := make([]float64, 0, p.Paths)

	for i := 0; i < p.Paths; i++ {
		for j := range zs {
			zs[j] = n.Next()
		}
		path, err := rng.Path(p.Spot, p.Rate, p.Vol, dt, zs)
		if err != nil {
			return engine.Price{}, err
		}
		touched := barrierTouched(path, p.Barrier, p.BarrierType)
		payoff := computePayoff(path, p.Strike, p.IsCall, touched, p.BarrierType)
		payoffs = append(payoffs, disc*payoff)
	}

	mean, sd := sampleStats(payoffs)
	return engine.Price{Value: mean, StdErr: sd / math.Sqrt(float64(len(payoffs)))}, nil
}

func barrierTouched(path []float64, barrier float64, bt BarrierType) bool {
	for _, s := range path {
		switch bt {
		case UpAndOut, UpAndIn:
			if s >= barrier {
				return true
			}
		case DownAndOut, DownAndIn:
			if s <= barrier {
				return true
			}
		}
	}
	return false
}

func computePayoff(path []float64, strike float64, isCall, touched bool, bt BarrierType) float64 {
	final := path[len(path)-1]
	var intrinsic float64
	if isCall {
		intrinsic = math.Max(final-strike, 0)
	} else {
		intrinsic = math.Max(strike-final, 0)
	}
	switch bt {
	case UpAndOut, DownAndOut:
		if touched {
			return 0
		}
		return intrinsic
	case UpAndIn, DownAndIn:
		if touched {
			return intrinsic
		}
		return 0
	}
	return 0
}

func sampleStats(xs []float64) (mean, sd float64) {
	for _, x := range xs {
		mean += x
	}
	mean /= float64(len(xs))
	for _, x := range xs {
		sd += (x - mean) * (x - mean)
	}
	if len(xs) > 1 {
		sd = math.Sqrt(sd / float64(len(xs)-1))
	}
	return mean, sd
}
