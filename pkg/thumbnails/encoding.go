package thumbnails

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

type Encoder interface {
	Encode(io.Writer, image.Image) error
	Types() []string
}

type PngEncoder struct{}

func (e PngEncoder) Encode(w io.Writer, i image.Image) error {
	return png.Encode(w, i)
}

func (e PngEncoder) Types() []string {
	return []string{"png"}
}

type JpegEncoder struct{}

func (e JpegEncoder) Encode(w io.Writer, i image.Image) error {
	return jpeg.Encode(w, i, nil)
}

func (e JpegEncoder) Types() []string {
	return []string{"jpeg", "jpg"}
}

func EncoderForType(fileType string) Encoder {
	switch strings.ToLower(fileType) {
	case "png":
		return PngEncoder{}
	case "jpg":
		fallthrough
	case "jpeg":
		return JpegEncoder{}
	default:
		return nil
	}
}
