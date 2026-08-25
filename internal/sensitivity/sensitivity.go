package sensitivity

import (
	"errors"
	"math"

	"mc-option/internal/engine"
)

var ErrInvalidGrid = errors.New("sensitivity: invalid grid parameters")

type GridPoint struct {
	ParamValue float64 `json:"param_value"`
	Price      float64 `json:"price"`
	StdErr     float64 `json:"stderr"`
}

type SweepResult struct {
	ParamName string      `json:"param_name"`
	Points    []GridPoint `json:"points"`
	MinPrice  float64     `json:"min_price"`
	MaxPrice  float64     `json:"max_price"`
}

func SpotSweep(base engine.Params, isCall bool, lo, hi float64, steps int) (*SweepResult, error) {
	return sweep(base, isCall, "spot", lo, hi, steps, func(p *engine.Params, v float64) { p.Spot = v })
}

func VolSweep(base engine.Params, isCall bool, lo, hi float64, steps int) (*SweepResult, error) {
	return sweep(base, isCall, "vol", lo, hi, steps, func(p *engine.Params, v float64) { p.Vol = v })
}

func StrikeSweep(base engine.Params, isCall bool, lo, hi float64, steps int) (*SweepResult, error) {
	return sweep(base, isCall, "strike", lo, hi, steps, func(p *engine.Params, v float64) { p.Strike = v })
}

func MaturitySweep(base engine.Params, isCall bool, lo, hi float64, steps int) (*SweepResult, error) {
	return sweep(base, isCall, "maturity", lo, hi, steps, func(p *engine.Params, v float64) { p.Maturity = v })
}

func RateSweep(base engine.Params, isCall bool, lo, hi float64, steps int) (*SweepResult, error) {
	return sweep(base, isCall, "rate", lo, hi, steps, func(p *engine.Params, v float64) { p.Rate = v })
}

func sweep(base engine.Params, isCall bool, name string, lo, hi float64, steps int, set func(*engine.Params, float64)) (*SweepResult, error) {
	if steps < 2 || lo >= hi {
		return nil, ErrInvalidGrid
	}
	res := &SweepResult{ParamName: name, MinPrice: math.MaxFloat64}
	step := (hi - lo) / float64(steps-1)
	var slots [][]float64
	var params []float64
	var stderrs []float64
	for i := 0; i < steps; i++ {
		v := lo + float64(i)*step
		p := base
		set(&p, v)
		pr, err := engine.European(p, isCall)
		if err != nil {
			continue
		}
		slots = append(slots, engine.PublishSweepPrice(pr.Value))
		params = append(params, v)
		stderrs = append(stderrs, pr.StdErr)
		if pr.Value < res.MinPrice {
			res.MinPrice = pr.Value
		}
		if pr.Value > res.MaxPrice {
			res.MaxPrice = pr.Value
		}
	}
	for i := range slots {
		res.Points = append(res.Points, GridPoint{
			ParamValue: params[i],
			Price:      slots[i][0],
			StdErr:     stderrs[i],
		})
	}
	return res, nil
}

type SurfacePoint struct {
	Strike   float64 `json:"strike"`
	Maturity float64 `json:"maturity"`
	Price    float64 `json:"price"`
}

func PriceSurface(base engine.Params, isCall bool, strikes []float64, maturities []float64) ([]SurfacePoint, error) {
	var points []SurfacePoint
	for _, k := range strikes {
		for _, t := range maturities {
			p := base
			p.Strike = k
			p.Maturity = t
			pr, err := engine.European(p, isCall)
			if err != nil {
				continue
			}
			points = append(points, SurfacePoint{Strike: k, Maturity: t, Price: pr.Value})
		}
	}
	return points, nil
}
