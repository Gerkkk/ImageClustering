package domain

import "math"

type ColorDistance struct{}

func (cd *ColorDistance) CIE76(p1, p2 Pixel) float64 {
	L1, a1, b1 := rgbToLab(p1.R, p1.G, p1.B)
	L2, a2, b2 := rgbToLab(p2.R, p2.G, p2.B)

	return math.Sqrt(math.Pow(L1-L2, 2) + math.Pow(a1-a2, 2) + math.Pow(b1-b2, 2))
}

func rgbToLab(r, g, b uint8) (float64, float64, float64) {
	R := float64(r) / 255.0
	G := float64(g) / 255.0
	B := float64(b) / 255.0

	R = gammaCorrect(R)
	G = gammaCorrect(G)
	B = gammaCorrect(B)

	X := 0.4124*R + 0.3576*G + 0.1805*B
	Y := 0.2126*R + 0.7152*G + 0.0722*B
	Z := 0.0193*R + 0.1192*G + 0.9505*B

	Xn, Yn, Zn := 0.9505, 1.0, 1.089
	L := 116*f(Y/Yn) - 16
	a := 500 * (f(X/Xn) - f(Y/Yn))
	br := 200 * (f(Y/Yn) - f(Z/Zn))
	return L, a, br
}

func f(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3.0)
	}
	return (t / 903.3) + 16.0/116.0
}

func gammaCorrect(c float64) float64 {
	if c > 0.04045 {
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	return c / 12.92
}
