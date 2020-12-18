package thumbnail

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

// Encoder encodes the thumbnail to a specific format.
type Encoder interface {
	// Encode encodes the image to a format.
	Encode(io.Writer, image.Image) error
	// Types returns the formats suffixes.
	Types() []string
	// MimeType returns the mimetype used by the encoder.
	MimeType() string
}

// PngEncoder encodes to png
type PngEncoder struct{}

// Encode encodes to png format
func (e PngEncoder) Encode(w io.Writer, i image.Image) error {
	return png.Encode(w, i)
}

// Types returns the png suffix
func (e PngEncoder) Types() []string {
	return []string{"png"}
}

// MimeType returns the mimetype for png files.
func (e PngEncoder) MimeType() string {
	return "image/png"
}

// JpegEncoder encodes to jpg.
type JpegEncoder struct{}

// Encode encodes to jpg
func (e JpegEncoder) Encode(w io.Writer, i image.Image) error {
	return jpeg.Encode(w, i, nil)
}

// Types returns the jpg suffixes.
func (e JpegEncoder) Types() []string {
	return []string{"jpeg", "jpg"}
}

// MimeType returns the mimetype for jpg files.
func (e JpegEncoder) MimeType() string {
	return "image/jpeg"
}

// EncoderForType returns the encoder for a given file type
// or nil if the type is not supported.
func EncoderForType(fileType string) Encoder {
	switch strings.ToLower(fileType) {
	case "png":
		return PngEncoder{}
	case "jpg", "jpeg":
		return JpegEncoder{}
	default:
		return nil
	}
}
