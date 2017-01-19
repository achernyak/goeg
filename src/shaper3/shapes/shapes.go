package shapes

import "image/color"

func clamp(minimum, x, maximum int) int {
	switch {
	case x < minimum:
		return minimum
	case x > maximum:
		return maximum
	}
	return x
}

func validFillColor(fill color.Color) color.Color {
	if fill == nil {
		return color.Black
	}
	return fill
}
