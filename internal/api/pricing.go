package api

import "mc-option/internal/engine"

type BatchPriceRequest struct {
	Requests []PriceRequest `json:"requests"`
}

type BatchPriceResponse struct {
	Results []PriceResponse `json:"results"`
	Errors  []string        `json:"errors,omitempty"`
}

func ValidateParams(req PriceRequest) (*engine.Params, error) {
	p := engine.Params{
		Spot:     req.Spot,
		Vol:      req.Vol,
		Rate:     req.Rate,
		Strike:   req.Strike,
		Maturity: req.Maturity,
		Steps:    req.Steps,
		Paths:    req.Paths,
		Seed:     req.Seed,
	}
	if p.Steps == 0 {
		p.Steps = 64
	}
	if p.Paths == 0 {
		p.Paths = 10000
	}
	if p.Seed == 0 {
		p.Seed = 42
	}
	if err := engine.Validate(p); err != nil {
		return nil, err
	}
	return &p, nil
}

func OptionTypeIsCall(typ string) bool {
	return typ == "euro-call" || typ == "asian-call"
}

func OptionTypeIsAsian(typ string) bool {
	return typ == "asian-call" || typ == "asian-put"
}

func SupportedTypes() []string {
	return []string{"euro-call", "euro-put", "asian-call", "asian-put"}
}

func IsValidType(typ string) bool {
	for _, t := range SupportedTypes() {
		if t == typ {
			return true
		}
	}
	return false
}
