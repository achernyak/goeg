package main

import (
	"image"
	"image/color"
	"shaper3/shapes"
)

func main() {
	img := shapes.FilledImage(420, 220, image.White)
	fill := color.RGBA{200, 200, 200, 0xFF}
	for i := 0; i < 10; i++ {
		width, height := 40+(20*i), 20+(10*i)
		rectangle := shapes.Rectangle{fill,
			image.Rect(0, 0, width, height), true}
		x := 10 + (20 * i)
		for j := i / 2; j >= 0; j-- {
			rectangle.Draw(img, x+j, (x/2)+j)
		}
		fill.R -= uint8(i * 5)
		fill.G = fill.R
		fill.B = fill.R
	}
	shapes.SaveImage(img, "rectangle.png")
}
