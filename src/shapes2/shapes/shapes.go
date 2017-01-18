package shapes

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
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

func checkBounds(img image.Image, x, y int) error {
	if !image.Rect(x, y, x, y).In(img.Bounds()) {
		return fmt.Errorf("%s(): point (%d, %d) is outside the image\n",
			caller(1), x, y)
	}
	return nil
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
