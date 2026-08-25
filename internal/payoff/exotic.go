package payoff

import "math"

type BestOf struct {
	Strike float64
}

func (b BestOf) Compute(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	return math.Max(prices[len(prices)-1]-b.Strike, 0)
}

func (b BestOf) Name() string { return "best-of" }

type Chooser struct {
	Strike float64
}

func (c Chooser) Compute(prices []float64) float64 {
	final := prices[len(prices)-1]
	callPay := math.Max(final-c.Strike, 0)
	putPay := math.Max(c.Strike-final, 0)
	return math.Max(callPay, putPay)
}

func (c Chooser) Name() string { return "chooser" }

type PowerCall struct {
	Strike float64
	Power  float64
}

func (p PowerCall) Compute(prices []float64) float64 {
	final := prices[len(prices)-1]
	return math.Max(math.Pow(final, p.Power)-p.Strike, 0)
}

func (p PowerCall) Name() string { return "power-call" }

type Cliquet struct {
	Floor float64
	Cap   float64
}

func (cl Cliquet) Compute(prices []float64) float64 {
	if len(prices) < 2 {
		return 0
	}
	total := 0.0
	for i := 1; i < len(prices); i++ {
		ret := (prices[i] - prices[i-1]) / prices[i-1]
		if ret < cl.Floor {
			ret = cl.Floor
		}
		if ret > cl.Cap {
			ret = cl.Cap
		}
		total += ret
	}
	return math.Max(total, 0) * prices[0]
}

func (cl Cliquet) Name() string { return "cliquet" }

type Spread struct {
	StrikeLow  float64
	StrikeHigh float64
}

func (s Spread) Compute(prices []float64) float64 {
	final := prices[len(prices)-1]
	return math.Max(final-s.StrikeLow, 0) - math.Max(final-s.StrikeHigh, 0)
}

func (s Spread) Name() string { return "call-spread" }

type Butterfly struct {
	StrikeLow  float64
	StrikeMid  float64
	StrikeHigh float64
}

func (bf Butterfly) Compute(prices []float64) float64 {
	final := prices[len(prices)-1]
	return math.Max(final-bf.StrikeLow, 0) - 2*math.Max(final-bf.StrikeMid, 0) + math.Max(final-bf.StrikeHigh, 0)
}

func (bf Butterfly) Name() string { return "butterfly" }
