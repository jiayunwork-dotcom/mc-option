package risk

import (
	"math"
	"sort"
)

type Position struct {
	Name     string
	Quantity float64
	PnL      []float64
}

type Portfolio struct {
	Positions []Position
}

func NewPortfolio() *Portfolio {
	return &Portfolio{}
}

func (pf *Portfolio) Add(pos Position) {
	pf.Positions = append(pf.Positions, pos)
}

func (pf *Portfolio) CombinedPnL() []float64 {
	if len(pf.Positions) == 0 {
		return nil
	}
	n := len(pf.Positions[0].PnL)
	combined := make([]float64, n)
	for _, pos := range pf.Positions {
		for i := 0; i < n && i < len(pos.PnL); i++ {
			combined[i] += pos.Quantity * pos.PnL[i]
		}
	}
	return combined
}

func (pf *Portfolio) PortfolioVaR(conf float64) (float64, error) {
	return VaR(pf.CombinedPnL(), conf)
}

func (pf *Portfolio) PortfolioES(conf float64) (float64, error) {
	return ES(pf.CombinedPnL(), conf)
}

func (pf *Portfolio) DiversificationBenefit(conf float64) (float64, error) {
	sumIndividual := 0.0
	for _, pos := range pf.Positions {
		scaled := make([]float64, len(pos.PnL))
		for i, v := range pos.PnL {
			scaled[i] = pos.Quantity * v
		}
		v, err := VaR(scaled, conf)
		if err != nil {
			return 0, err
		}
		sumIndividual += v
	}
	combined, err := pf.PortfolioVaR(conf)
	if err != nil {
		return 0, err
	}
	return sumIndividual - combined, nil
}

type Scenario struct {
	Name      string
	SpotShift float64
	VolShift  float64
}

type ScenarioResult struct {
	ScenarioName string  `json:"scenario"`
	PnLChange    float64 `json:"pnl_change"`
}

func Percentiles(xs []float64, pcts []float64) []float64 {
	sorted := append([]float64(nil), xs...)
	sort.Float64s(sorted)
	out := make([]float64, len(pcts))
	n := float64(len(sorted))
	for i, p := range pcts {
		idx := p / 100 * (n - 1)
		lo := int(math.Floor(idx))
		hi := int(math.Ceil(idx))
		if lo < 0 {
			lo = 0
		}
		if hi >= len(sorted) {
			hi = len(sorted) - 1
		}
		frac := idx - float64(lo)
		out[i] = sorted[lo]*(1-frac) + sorted[hi]*frac
	}
	return out
}

func MaxDrawdown(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	peak := xs[0]
	maxDD := 0.0
	for _, x := range xs {
		if x > peak {
			peak = x
		}
		dd := peak - x
		if dd > maxDD {
			maxDD = dd
		}
	}
	return maxDD
}

func SharpeRatio(pnl []float64, riskFreeReturn float64) float64 {
	if len(pnl) < 2 {
		return 0
	}
	sum := 0.0
	for _, x := range pnl {
		sum += x - riskFreeReturn
	}
	mean := sum / float64(len(pnl))
	sumSq := 0.0
	for _, x := range pnl {
		d := (x - riskFreeReturn) - mean
		sumSq += d * d
	}
	std := math.Sqrt(sumSq / float64(len(pnl)-1))
	if std == 0 {
		return 0
	}
	return mean / std
}
