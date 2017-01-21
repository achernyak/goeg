package font

import (
	"fmt"
	"log"
	"unicode/utf8"
)

type Font struct {
	family string
	size   int
}

func New(family string, size int) *Font {
	return &Font{saneFamily("sans-serif", family), saneSize(10, size)}
}

func (font *Font) Family() string { return font.family }

func (font *Font) SetFamily(family string) {
	font.family = saneFamily(font.family, family)
}

func (font *Font) Size() int { return font.size }

func (font *Font) SetSize(size int) {
	font.size = saneSize(font.size, size)
}

func (font *Font) String() string {
	return fmt.Sprintf("{font-family: %q; font-size: %dpt;}", font.family,
		font.size)
}

func saneFamily(oldFamily, newFamily string) string {
	if len(newFamily) < utf8.UTFMax &&
		utf8.RuneCountInString(newFamily) < 1 {
		log.Printf("font.saneFamily(): ignored invalid family '%s'",
			newFamily)
		return oldFamily
	}
	return newFamily
}

func saneSize(oldSize, newSize int) int {
	if newSize < 5 || newSize > 144 {
		log.Printf("font.saneSize(): ignored invalid size '%d'", newSize)
		return oldSize
	}
	return newSize
}
