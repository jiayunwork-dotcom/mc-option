package rng

import "math"

func Halton(n int, base int) float64 {
	f := 1.0
	result := 0.0
	i := n
	for i > 0 {
		f /= float64(base)
		result += f * float64(i%base)
		i /= base
	}
	return result
}

func HaltonSequence(count, base int) []float64 {
	out := make([]float64, count)
	for i := 0; i < count; i++ {
		out[i] = Halton(i+1, base)
	}
	return out
}

func HaltonNormal(count, base int) []float64 {
	uniform := HaltonSequence(count, base)
	out := make([]float64, count)
	for i, u := range uniform {
		out[i] = InvNormCDF(u)
	}
	return out
}

func InvNormCDF(u float64) float64 {
	if u <= 0 {
		return -10
	}
	if u >= 1 {
		return 10
	}
	const (
		a1 = -3.969683028665376e+01
		a2 = 2.209460984245205e+02
		a3 = -2.759285104469687e+02
		a4 = 1.383577518672690e+02
		a5 = -3.066479806614716e+01
		a6 = 2.506628277459239e+00
		b1 = -5.447609879822406e+01
		b2 = 1.615858368580409e+02
		b3 = -1.556989798598866e+02
		b4 = 6.680131188771972e+01
		b5 = -1.328068155288572e+01
		c1 = -7.784894002430293e-03
		c2 = -3.223964580411365e-01
		c3 = -2.400758277161838e+00
		c4 = -2.549732539343734e+00
		c5 = 4.374664141464968e+00
		c6 = 2.938163982698783e+00
		d1 = 7.784695709041462e-03
		d2 = 3.224671290700398e-01
		d3 = 2.445134137142996e+00
		d4 = 3.754408661907416e+00
	)
	pLow := 0.02425
	pHigh := 1 - pLow
	var q, r float64
	if u < pLow {
		q = math.Sqrt(-2 * math.Log(u))
		return (((((c1*q+c2)*q+c3)*q+c4)*q+c5)*q + c6) / ((((d1*q+d2)*q+d3)*q+d4)*q + 1)
	}
	if u <= pHigh {
		q = u - 0.5
		r = q * q
		return (((((a1*r+a2)*r+a3)*r+a4)*r+a5)*r + a6) * q / (((((b1*r+b2)*r+b3)*r+b4)*r+b5)*r + 1)
	}
	q = math.Sqrt(-2 * math.Log(1-u))
	return -(((((c1*q+c2)*q+c3)*q+c4)*q+c5)*q + c6) / ((((d1*q+d2)*q+d3)*q+d4)*q + 1)
}

func VanDerCorput(n, base int) float64 {
	return Halton(n, base)
}

func StratifiedSampling(n int, src *Normal) []float64 {
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		u := (float64(i) + src.src.Float64()) / float64(n)
		out[i] = InvNormCDF(u)
	}
	return out
}
