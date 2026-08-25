package bs

import (
	"errors"
	"math"
)

var ErrNoConverge = errors.New("bs: implied vol did not converge")

func ImpliedVol(spot, strike, rate, maturity, marketPrice float64, isCall bool) (float64, error) {
	if marketPrice <= 0 || spot <= 0 || strike <= 0 || maturity <= 0 {
		return 0, ErrInvalidInput
	}
	lo, hi := 0.001, 5.0
	for iter := 0; iter < 200; iter++ {
		mid := (lo + hi) / 2
		var pr float64
		if isCall {
			r, err := Call(spot, strike, rate, mid, maturity)
			if err != nil {
				return 0, err
			}
			pr = r.Price
		} else {
			r, err := Put(spot, strike, rate, mid, maturity)
			if err != nil {
				return 0, err
			}
			pr = r.Price
		}
		if math.Abs(pr-marketPrice) < 1e-8 {
			return mid, nil
		}
		if pr < marketPrice {
			lo = mid
		} else {
			hi = mid
		}
		if hi-lo < 1e-10 {
			break
		}
	}
	return (lo + hi) / 2, nil
}

func ImpliedVolNewton(spot, strike, rate, maturity, marketPrice float64, isCall bool) (float64, error) {
	if marketPrice <= 0 || spot <= 0 || strike <= 0 || maturity <= 0 {
		return 0, ErrInvalidInput
	}
	vol := 0.2
	for iter := 0; iter < 100; iter++ {
		var pr float64
		if isCall {
			r, _ := Call(spot, strike, rate, vol, maturity)
			pr = r.Price
		} else {
			r, _ := Put(spot, strike, rate, vol, maturity)
			pr = r.Price
		}
		v, _ := Vega(spot, strike, rate, vol, maturity)
		vegaFull := v * 100
		if math.Abs(vegaFull) < 1e-12 {
			return vol, ErrNoConverge
		}
		diff := pr - marketPrice
		if math.Abs(diff) < 1e-8 {
			return vol, nil
		}
		vol -= diff / vegaFull
		if vol <= 0 {
			vol = 0.001
		}
	}
	return vol, ErrNoConverge
}

type SmilePoint struct {
	Strike float64 `json:"strike"`
	IV     float64 `json:"iv"`
}

func ComputeSmile(spot, rate, maturity float64, strikes, marketPrices []float64, isCall bool) ([]SmilePoint, error) {
	if len(strikes) != len(marketPrices) {
		return nil, errors.New("bs: strikes and prices must have same length")
	}
	var points []SmilePoint
	for i, k := range strikes {
		iv, err := ImpliedVol(spot, k, rate, maturity, marketPrices[i], isCall)
		if err != nil {
			continue
		}
		points = append(points, SmilePoint{Strike: k, IV: iv})
	}
	return points, nil
}

func Moneyness(spot, strike float64) float64 {
	if strike == 0 {
		return 0
	}
	return spot / strike
}

func LogMoneyness(spot, strike float64) float64 {
	if strike <= 0 || spot <= 0 {
		return 0
	}
	return math.Log(spot / strike)
}
