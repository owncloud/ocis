package thumbnail

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"strings"

	"github.com/disintegration/imaging"
)

var (
	// ErrInvalidType represents the error when a type can't be encoded.
	ErrInvalidType2 = errors.New("can't encode this type")
	// ErrNoGeneratorForType represents the error when no generator could be found for a type.
	ErrNoGeneratorForType = errors.New("no generator for this type found")
)

type Generator interface {
	GenerateThumbnail(image.Rectangle, interface{}) (interface{}, error)
}

type SimpleGenerator struct{}

func (g SimpleGenerator) GenerateThumbnail(size image.Rectangle, img interface{}) (interface{}, error) {
	m, ok := img.(image.Image)
	if !ok {
		return nil, ErrInvalidType2
	}

	return imaging.Thumbnail(m, size.Dx(), size.Dy(), imaging.Lanczos), nil
}

type GifGenerator struct{}

func (g GifGenerator) GenerateThumbnail(size image.Rectangle, img interface{}) (interface{}, error) {
	// Code inspired by https://github.com/willnorris/gifresize/blob/db93a7e1dcb1c279f7eeb99cc6d90b9e2e23e871/gifresize.go

	m, ok := img.(*gif.GIF)
	if !ok {
		return nil, ErrInvalidType2
	}
	// Create a new RGBA image to hold the incremental frames.
	srcX, srcY := m.Config.Width, m.Config.Height
	b := image.Rect(0, 0, srcX, srcY)
	tmp := image.NewRGBA(b)

	for i, frame := range m.Image {
		bounds := frame.Bounds()
		prev := tmp
		draw.Draw(tmp, bounds, frame, bounds.Min, draw.Over)
		scaled := imaging.Resize(tmp, size.Dx(), size.Dy(), imaging.Lanczos)
		m.Image[i] = g.imageToPaletted(scaled, frame.Palette)

		switch m.Disposal[i] {
		case gif.DisposalBackground:
			tmp = image.NewRGBA(b)
		case gif.DisposalPrevious:
			tmp = prev
		}
	}
	m.Config.Width = size.Dx()
	m.Config.Height = size.Dy()

	return m, nil
}

func (g GifGenerator) imageToPaletted(img image.Image, p color.Palette) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, p)
	draw.FloydSteinberg.Draw(pm, b, img, image.Point{})
	return pm
}

// GeneratorForType returns the generator for a given file type
// or nil if the type is not supported.
func GeneratorForType(fileType string) (Generator, error) {
	switch strings.ToLower(fileType) {
	case typePng, typeJpg, typeJpeg:
		return SimpleGenerator{}, nil
	case typeGif:
		return GifGenerator{}, nil
	default:
		return nil, ErrNoEncoderForType
	}
}
