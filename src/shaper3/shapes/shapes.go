package shapes

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
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

type RegularPolygon struct {
	color.Color
	Radius int
	Sides  int
}

func (polygon RegularPolygon) Draw(img draw.Image, x, y int) error {
	if err := checkBounds(img, x, y); err != nil {
		return err

	}
	fill := validFillColor(polygon.Color)
	radius := clamp(1, polygon.Radius, 1024)
	sides := clamp(3, polygon.Sides, 60)
	points := getPoints(x, y, sides, float64(radius))
	for i := 0; i < sides; i++ { // Draw lines between the apexes
		drawLine(img, points[i], points[i+1], fill)
	}
	return nil
}

func getPoints(x, y, sides int, radius float64) []image.Point {
	points := make([]image.Point, sides+1)
	// Compute the shape's apexes (thanks to Jasmin Blanchette)
	fullCircle := 2 * math.Pi
	x0, y0 := float64(x), float64(y)
	for i := 0; i < sides; i++ {
		θ := float64(float64(i) * fullCircle / float64(sides))
		x1 := x0 + (radius * math.Sin(θ))
		y1 := y0 + (radius * math.Cos(θ))
		points[i] = image.Pt(int(x1), int(y1))
	}
	points[sides] = points[0] // close the shape
	return points
}

func drawLine(img draw.Image, start, end image.Point,
	fill color.Color) {
	x0, x1 := start.X, end.X
	y0, y1 := start.Y, end.Y
	Δx := math.Abs(float64(x1 - x0))
	Δy := math.Abs(float64(y1 - y0))
	if Δx >= Δy { // shallow slope
		if x0 > x1 {
			x0, y0, x1, y1 = x1, y1, x0, y0

		}
		y := y0
		yStep := 1
		if y0 > y1 {
			yStep = -1

		}
		remainder := float64(int(Δx/2)) - Δx
		for x := x0; x <= x1; x++ {
			img.Set(x, y, fill)
			remainder += Δy
			if remainder >= 0.0 {
				remainder -= Δx
				y += yStep

			}

		}
	} else { // steep slope
		if y0 > y1 {
			x0, y0, x1, y1 = x1, y1, x0, y0

		}
		x := x0
		xStep := 1
		if x0 > x1 {
			xStep = -1

		}
		remainder := float64(int(Δy/2)) - Δy
		for y := y0; y <= y1; y++ {
			img.Set(x, y, fill)
			remainder += Δx
			if remainder >= 0.0 {
				remainder -= Δy
				x += xStep

			}

		}
	}

}

func (polygon RegularPolygon) String() string {
	return fmt.Sprintf("polygon(fill=%v, radius=%d, sides=%d)",
		polygon.Color, polygon.Radius, polygon.Sides)
}
