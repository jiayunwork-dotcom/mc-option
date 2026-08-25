package payoff

import "math"

type Payoffer interface {
	Compute(prices []float64) float64
	Name() string
}

type VanillaCall struct {
	Strike float64
}

func (v VanillaCall) Compute(prices []float64) float64 {
	return math.Max(prices[len(prices)-1]-v.Strike, 0)
}

func (v VanillaCall) Name() string { return "vanilla-call" }

type VanillaPut struct {
	Strike float64
}

func (v VanillaPut) Compute(prices []float64) float64 {
	return math.Max(v.Strike-prices[len(prices)-1], 0)
}

func (v VanillaPut) Name() string { return "vanilla-put" }

type AsianCall struct {
	Strike float64
}

func (a AsianCall) Compute(prices []float64) float64 {
	avg := arithmeticMean(prices)
	return math.Max(avg-a.Strike, 0)
}

func (a AsianCall) Name() string { return "asian-call" }

type AsianPut struct {
	Strike float64
}

func (a AsianPut) Compute(prices []float64) float64 {
	avg := arithmeticMean(prices)
	return math.Max(a.Strike-avg, 0)
}

func (a AsianPut) Name() string { return "asian-put" }

type LookbackCall struct{}

func (l LookbackCall) Compute(prices []float64) float64 {
	min := prices[0]
	for _, s := range prices[1:] {
		if s < min {
			min = s
		}
	}
	return math.Max(prices[len(prices)-1]-min, 0)
}

func (l LookbackCall) Name() string { return "lookback-call" }

type LookbackPut struct{}

func (l LookbackPut) Compute(prices []float64) float64 {
	max := prices[0]
	for _, s := range prices[1:] {
		if s > max {
			max = s
		}
	}
	return math.Max(max-prices[len(prices)-1], 0)
}

func (l LookbackPut) Name() string { return "lookback-put" }

type DigitalCall struct {
	Strike float64
	Amount float64
}

func (d DigitalCall) Compute(prices []float64) float64 {
	if prices[len(prices)-1] > d.Strike {
		return d.Amount
	}
	return 0
}

func (d DigitalCall) Name() string { return "digital-call" }

type DigitalPut struct {
	Strike float64
	Amount float64
}

func (d DigitalPut) Compute(prices []float64) float64 {
	if prices[len(prices)-1] < d.Strike {
		return d.Amount
	}
	return 0
}

func (d DigitalPut) Name() string { return "digital-put" }

type Straddle struct {
	Strike float64
}

func (s Straddle) Compute(prices []float64) float64 {
	final := prices[len(prices)-1]
	return math.Abs(final - s.Strike)
}

func (s Straddle) Name() string { return "straddle" }

func arithmeticMean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}
