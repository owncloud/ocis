package thumbnail

import (
	"image"
	"strings"

	"github.com/kovidgoyal/imaging"
)

// Processor processes the thumbnail by applying different transformations to it.
type Processor interface {
	ID() string
	Process(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA
}

// DefinableProcessor is the simplest processor, it holds a replaceable image converter function.
type DefinableProcessor struct {
	Slug      string
	Converter func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA
}

// ID returns the processor identification.
func (p DefinableProcessor) ID() string { return p.Slug }

// Process transforms the given image.
func (p DefinableProcessor) Process(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA {
	return p.Converter(img, width, height, filter)
}

// ProcessorFor returns a matching Processor
func ProcessorFor(id, fileType string) (DefinableProcessor, error) {
	convertToNRGBA := func(img image.Image) *image.NRGBA {
		if nrgba, ok := img.(*image.NRGBA); ok {
			return nrgba
		}
		return imaging.Clone(img)
	}
	switch strings.ToLower(id) {
	case "fit":
		return DefinableProcessor{Slug: strings.ToLower(id), Converter: func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA {
			return convertToNRGBA(imaging.Fit(img, width, height, filter))
		}}, nil
	case "resize":
		return DefinableProcessor{Slug: strings.ToLower(id), Converter: func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA {
			return convertToNRGBA(imaging.Resize(img, width, height, filter))
		}}, nil
	case "fill":
		return DefinableProcessor{Slug: strings.ToLower(id), Converter: func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA {
			return convertToNRGBA(imaging.Fill(img, width, height, imaging.Center, filter))
		}}, nil
	case "thumbnail":
		return DefinableProcessor{Slug: strings.ToLower(id), Converter: func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA {
			return convertToNRGBA(imaging.Thumbnail(img, width, height, filter))
		}}, nil
	default:
		switch strings.ToLower(fileType) {
		case typeGif:
			return DefinableProcessor{Converter: func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA {
				return convertToNRGBA(imaging.Resize(img, width, height, filter))
			}}, nil
		default:
			return DefinableProcessor{Converter: func(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA {
				return convertToNRGBA(imaging.Thumbnail(img, width, height, filter))
			}}, nil
		}
	}
}
