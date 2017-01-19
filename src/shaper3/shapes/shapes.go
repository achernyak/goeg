package shapes

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"path/filepath"
	"runtime"
)

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

type Drawer interface {
	Draw(img draw.Image, x, y int) error
}

type Circle struct {
	color.Color
	Radius int
}

func (circle Circle) Draw(img draw.Image, x, y int) error {
	if err := checkBounds(img, x, y); err != nil {
		return err

	}
	fill := validFillColor(circle.Color)
	radius := clamp(1, circle.Radius, 1024)

	x0, y0 := x, y
	f := 1 - radius
	ddF_x, ddF_y := 1, -2*radius
	x, y = 0, radius

	img.Set(x0, y0+radius, fill)
	img.Set(x0, y0-radius, fill)
	img.Set(x0+radius, y0, fill)
	img.Set(x0-radius, y0, fill)

	for x < y {
		if f >= 0 {
			y--
			ddF_y += 2
			f += ddF_y

		}
		x++
		ddF_x += 2
		f += ddF_x
		img.Set(x0+x, y0+y, fill)
		img.Set(x0-x, y0+y, fill)
		img.Set(x0+x, y0-y, fill)
		img.Set(x0-x, y0-y, fill)
		img.Set(x0+y, y0+x, fill)
		img.Set(x0-y, y0+x, fill)
		img.Set(x0+y, y0-x, fill)
		img.Set(x0-y, y0-x, fill)

	}
	return nil
}

func (circle Circle) String() string {
	return fmt.Sprintf("circle(fill=%v, radius=%d)", circle.Color,
		circle.Radius)
}

func checkBounds(img image.Image, x, y int) error {
	if !image.Rect(x, y, x, y).In(img.Bounds()) {
		return fmt.Errorf("%s(): point(%d, %d) is outside the image\n",
			caller(1), x, y)
	}
	return nil
}

func caller(steps int) string {
	name := "?"
	if pc, _, _, ok := runtime.Caller(steps + 1); ok {
		name = filepath.Base(runtime.FuncForPC(pc).Name())
	}
	return name
}
