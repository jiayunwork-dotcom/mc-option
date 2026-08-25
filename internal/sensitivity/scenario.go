package sensitivity

import "mc-option/internal/engine"

type Scenario struct {
	Name    string  `json:"name"`
	SpotPct float64 `json:"spot_pct"`
	VolAbs  float64 `json:"vol_abs"`
	RateAbs float64 `json:"rate_abs"`
}

type ScenarioResult struct {
	Scenario    Scenario `json:"scenario"`
	BasePrice   float64  `json:"base_price"`
	StressPrice float64  `json:"stress_price"`
	PnL         float64  `json:"pnl"`
}

func StandardScenarios() []Scenario {
	return []Scenario{
		{Name: "market-crash", SpotPct: -0.20, VolAbs: 0.15, RateAbs: -0.01},
		{Name: "vol-spike", SpotPct: -0.05, VolAbs: 0.20, RateAbs: 0},
		{Name: "rally", SpotPct: 0.15, VolAbs: -0.05, RateAbs: 0.005},
		{Name: "rate-hike", SpotPct: 0, VolAbs: 0, RateAbs: 0.02},
		{Name: "quiet-market", SpotPct: 0, VolAbs: -0.10, RateAbs: 0},
	}
}

func RunScenarios(base engine.Params, isCall bool, scenarios []Scenario) ([]ScenarioResult, error) {
	basePr, err := engine.European(base, isCall)
	if err != nil {
		return nil, err
	}
	var results []ScenarioResult
	for _, sc := range scenarios {
		p := base
		p.Spot *= (1 + sc.SpotPct)
		p.Vol += sc.VolAbs
		if p.Vol <= 0 {
			p.Vol = 0.01
		}
		p.Rate += sc.RateAbs
		if p.Rate <= 0 {
			p.Rate = 0.001
		}
		pr, err := engine.European(p, isCall)
		if err != nil {
			continue
		}
		results = append(results, ScenarioResult{
			Scenario:    sc,
			BasePrice:   basePr.Value,
			StressPrice: pr.Value,
			PnL:         pr.Value - basePr.Value,
		})
	}
	return results, nil
}

func WorstCase(results []ScenarioResult) *ScenarioResult {
	if len(results) == 0 {
		return nil
	}
	worst := &results[0]
	for i := 1; i < len(results); i++ {
		if results[i].PnL < worst.PnL {
			worst = &results[i]
		}
	}
	return worst
}

func BestCase(results []ScenarioResult) *ScenarioResult {
	if len(results) == 0 {
		return nil
	}
	best := &results[0]
	for i := 1; i < len(results); i++ {
		if results[i].PnL > best.PnL {
			best = &results[i]
		}
	}
	return best
}
