package blurhash

import "math"

func sRGBToLinear(value int) float64 {
	v := float64(value) / 255.0
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func linearToSRGB(value float64) int {
	v := math.Max(0, math.Min(1, value))
	if v <= 0.0031308 {
		return int(math.Round(v * 12.92 * 255))
	}
	return int(math.Round((1.055*math.Pow(v, 1.0/2.4) - 0.055) * 255))
}

func signPow(value, exp float64) float64 {
	if value < 0 {
		return -math.Pow(-value, exp)
	}
	return math.Pow(value, exp)
}

func encodeDC(r, g, b float64) int {
	return (linearToSRGB(r) << 16) + (linearToSRGB(g) << 8) + linearToSRGB(b)
}

func encodeAC(r, g, b, maximumValue float64) int {
	quantR := int(math.Max(0, math.Min(18, math.Floor(signPow(r/maximumValue, 0.5)*9+9.5))))
	quantG := int(math.Max(0, math.Min(18, math.Floor(signPow(g/maximumValue, 0.5)*9+9.5))))
	quantB := int(math.Max(0, math.Min(18, math.Floor(signPow(b/maximumValue, 0.5)*9+9.5))))
	return quantR*19*19 + quantG*19 + quantB
}
