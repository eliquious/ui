package ui

import (
	"io/ioutil"

	"github.com/flopp/go-findfont"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

// DefaultFontRegistry stores all the loaded fonts
var DefaultFontRegistry = NewFontRegistry()

// NewFontRegistry creates a new font face.
func NewFontRegistry() *FontRegistry {
	return &FontRegistry{
		make(map[string]*sfnt.Font),
	}
}

// FontRegistry is a registry of all the loaded fonts.
type FontRegistry struct {
	loadedFonts map[string]*sfnt.Font
}

// SystemFont loads a font into the registry from the OS.
func (r *FontRegistry) SystemFont(fontName string) (*sfnt.Font, error) {
	if f, ok := r.loadedFonts[fontName]; ok {
		return f, nil
	}

	// ffind font in system
	fontPath, err := findfont.Find(fontName)
	if err != nil {
		return nil, err
	}

	// load the font with the freetype library
	fontData, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	ff, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	r.loadedFonts[fontName] = ff
	return ff, nil
}

// NewFontFace creates a new font face from a system font.
func NewFontFace(fontName string, size float64) (font.Face, error) {
	tt, err := DefaultFontRegistry.SystemFont(fontName)
	if err != nil {
		return nil, err
	}

	return opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
