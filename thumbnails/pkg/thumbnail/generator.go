package thumbnail

import (
	"errors"
	"image"
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
	m, ok := img.(*gif.GIF)
	if !ok {
		return nil, ErrInvalidType2
	}
	var bounds image.Rectangle
	for i := range m.Image {
		img := imaging.Resize(m.Image[i], size.Dx(), size.Dy(), imaging.Lanczos)
		bounds = image.Rect(0, 0, size.Dx(), size.Dy())
		m.Image[i] = image.NewPaletted(bounds, m.Image[i].Palette)
		draw.Draw(m.Image[i], bounds, img, image.Pt(0, 0), draw.Src)
	}
	m.Config.Height = bounds.Dy()
	m.Config.Width = bounds.Dx()
	return m, nil
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
