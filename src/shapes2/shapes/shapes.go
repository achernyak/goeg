package shapes

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var saneLength, saneRadius, saneSides func(int) int

func init() {
	saneLength = makeBoundedIntFunc(1, 4096)
	saneRadius = makeBoundedIntFunc(1, 1024)
	saneSides = makeBoundedIntFunc(3, 60)
}

func makeBoundedIntFunc(minimum, maximum int) func(int) int {
	return func(x int) int {
		valid := x
		switch {
		case x < minimum:
			valid = minimum
		case x > maximum:
			valid = maximum
		}
		if valid != x {
			log.Printf("%s(): replace %d with %d\n", caller(1), x, valid)
		}
		return valid
	}
}

func caller(steps int) string {
	name := "?"
	if pc, _, _, ok := runtime.Caller(steps + 1); ok {
		name = filepath.Base(runtime.FuncForPC(pc).Name())
	}
	return name
}

type Shaper interface {
	Drawer
	Filler
}

type Drawer interface {
	Draw(img draw.Image, x, y int) error
}

type Filler interface {
	Fill() color.Color
	SetFill(fill color.Color)
}

type Radiuser interface {
	Radius() int
	SetRadius(radius int)
}

type Sideser interface {
	Sides() int
	SetSides(sides int)
}

type shape struct{ fill color.Color }

func newShape(fill color.Color) shape {
	if fill == nil {
		fill = color.Black
	}
	return shape{fill}
}

func (shape shape) Fill() color.Color { return shape.fill }

func (shape *shape) setFill(fill color.Color) {
	if fill == nil {
		fill = color.Black
	}
	shape.fill = fill
}

type Circle struct {
	shape
	radius int
}

func NewCircle(fill color.Color, radius int) *Circle {
	return &Circle{newShape(fill), saneRadius(radius)}
}

func (circle *Circle) Radius() int {
	return circle.radius
}

func (circle *Circle) SetRadius(radius int) {
	circle.radius = saneRadius(radius)
}

func (circle *Circle) Draw(img draw.Image, x, y int) error {
	if err := checkBounds(img, x, y); err != nil {
		return err
	}
	fill, radius := circle.fill, circle.radius
	x0, y0 := x, y
	f := 1 - radius
	ddF_x, ddF_y := 1, -2*radius
	x, y = 0, radius

	img.Set(x0, y0+radius, fill)
	img.Set(x0, y0-radius, fill)
	img.Set(y0+radius, x0, fill)
	img.Set(y0-radius, x0, fill)

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

func (circle *Circle) String() string {
	return fmt.Sprintf("circle(fill=%v, radius=%d)", circle.fill,
		circle.radius)
}

type RegularPolygon struct {
	*Circle
	sides int
}

func NewRegularPolygon(fill color.Color, radius,
	sides int) *RegularPolygon {
	return &RegularPolygon{NewCircle(fill, radius), saneSides(sides)}
}

func (polygon *RegularPolygon) Sides() int {
	return polygon.sides
}

func (polygon *RegularPolygon) SetSides(sides int) {
	polygon.sides = saneSides(sides)
}

func (polygon *RegularPolygon) Draw(img draw.Image, x, y int) error {
	if err := checkBounds(img, x, y); err != nil {
		return err
	}
	points := getPoints(x, y, polygon.sides, float64(polygon.Radius()))
	for i := 0; i < polygon.sides; i++ {
		drawLine(img, points[i], points[i+1], polygon.Fill())
	}
	return nil
}

func getPoints(x, y, sides int, radius float64) []image.Point {
	points := make([]image.Point, sides+1)
	fullCircle := 2 * math.Pi
	xangle, yangle := float64(x), float64(y)
	for i := 0; i < sides; i++ {
		angle := float64(float64(i) * fullCircle / float64(sides))
		x1 := xangle + (radius * math.Sin(angle))
		y1 := yangle + (radius * math.Cos(angle))
		points[i] = image.Pt(int(x1), int(y1))
	}
	points[sides] = points[0]
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

func (polygon *RegularPolygon) String() string {
	return fmt.Sprintf("polygon(fill=%v, radius=%d, sides=%d)",
		polygon.Fill(), polygon.Radius(), polygon.sides)
}

func checkBounds(img image.Image, x, y int) error {
	if !image.Rect(x, y, x, y).In(img.Bounds()) {
		return fmt.Errorf("%s(): point (%d, %d) is outside the image\n",
			caller(1), x, y)
	}
	return nil
}

type Option struct {
	Fill   color.Color
	Radius int
}

func New(shape string, option Option) (Shaper, error) {
	sidesForShape := map[string]int{"triangle": 3, "square": 4,
		"pentagon": 5, "hexagon": 6, "heptagon": 7, "octagon": 8,
		"enneagon": 9, "nonagon": 9, "decagon": 10}
	if sides, found := sidesForShape[shape]; found {
		return NewRegularPolygon(option.Fill, option.Radius, sides), nil
	}
	if shape != "circle" {
		return nil, fmt.Errorf("shapes.New(): invalide shape '%s'", shape)
	}
	return NewCircle(option.Fill, option.Radius), nil
}

func FillImage(width, height int, fill color.Color) draw.Image {
	if fill == nil {
		fill = color.Black
	}
	width = saneLength(width)
	height = saneLength(height)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{fill}, image.ZP, draw.Src)
	return img
}

func DrawShapes(img draw.Image, x, y int, shapes ...Shaper) error {
	for _, shape := range shapes {
		if err := shape.Draw(img, x, y); err != nil {
			return err
		}
	}
	return nil
}

func SaveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, nil)
	case ".png":
		return png.Encode(file, img)
	}
	return fmt.Errorf("shapes.SaveImage(): '%s' has an unrecognized "+
		"suffix", filename)
}
