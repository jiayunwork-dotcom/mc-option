package rng

import (
	"errors"
	"math"
	"math/rand"
)

type Normal struct {
	src      *rand.Rand
	spare    float64
	hasSpare bool
}

func NewNormal(seed int64) *Normal {
	return &Normal{src: rand.New(rand.NewSource(seed))}
}

func (n *Normal) Next() float64 {
	if n.hasSpare {
		n.hasSpare = false
		return n.spare
	}
	u1 := n.src.Float64()
	for u1 <= 0 {
		u1 = n.src.Float64()
	}
	u2 := n.src.Float64()
	r := math.Sqrt(-2 * math.Log(u1))
	theta := 2 * math.Pi * u2
	n.spare = r * math.Sin(theta)
	n.hasSpare = true
	return r * math.Cos(theta)
}

func (n *Normal) Antithetic(i int, z float64) float64 {
	return -z
}

func Path(spot, drift, vol, dt float64, zs []float64) ([]float64, error) {
	if len(zs) < 1 {
		return nil, errors.New("rng: zs must contain at least one draw")
	}
	path := make([]float64, len(zs)+1)
	path[0] = spot
	growth := growthFromLive(drift, vol, dt)
	diffusion := vol * math.Sqrt(dt)
	for i, z := range zs {
		path[i+1] = path[i] * growth * math.Exp(diffusion*z)
	}
	return path, nil
}
