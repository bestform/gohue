package gohue

import "math"

func convertRGBToXY(r, g, b int) (float64, float64) {
	red := float64(r) / 255.0
	green := float64(g) / 255.0
	blue := float64(b) / 255.0

	red = gammaCorrect(red)
	green = gammaCorrect(green)
	blue = gammaCorrect(blue)

	X := red*0.664511 + green*0.154324 + blue*0.162028
	Y := red*0.283881 + green*0.668433 + blue*0.047685
	Z := red*0.000088 + green*0.072310 + blue*0.986039

	x := X / (X + Y + Z)
	y := Y / (X + Y + Z)

	// @todo: Check if the found xy value is within the color gamut of the light
	// @todo: Use Y as brightness
	return x, y
}

func gammaCorrect(color float64) float64 {
	//float red = (red > 0.04045f) ? pow((red + 0.055f) / (1.0f + 0.055f), 2.4f) : (red / 12.92f);
	if color > 0.04045 {
		return math.Pow((color+0.055)/(1.0+0.055), 2.4)
	}

	return color / 12.92

}
