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

type Shape interface {
	Fill() color.Color
	SetFill(fill color.Color)
	Draw(img draw.Image, x, y int) error
}

type CircularShape interface {
	Shaper
	Radius() int
	SetRadius(radius int)
}

type RegularPolygonalShaper interface {
	CircularShaper
	Sides() int
	SetSides(sides int)
}

func caller(steps int) string {
	name := "?"
	if pc, _, _, ok := runtime.Caller(steps + 1); ok {
		name = filepath.Base(runtime.FuncForPC(pc).Name())
	}
	return name
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
