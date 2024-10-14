//go:build enable_vips

package thumbnail

import (
	"io"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/errors"
)

// PngEncoder encodes to png
type PngEncoder struct{}

// Encode encodes to png format
func (e PngEncoder) Encode(w io.Writer, img interface{}) error {
	m, ok := img.(*vips.ImageRef)
	if !ok {
		return errors.ErrInvalidType
	}

	buf, _, err := m.ExportPng(vips.NewPngExportParams())
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

// Types returns the png suffix
func (e PngEncoder) Types() []string {
	return []string{typePng}
}

// MimeType returns the mimetype for png files.
func (e PngEncoder) MimeType() string {
	return "image/png"
}

// JpegEncoder encodes to jpg
type JpegEncoder struct{}

// Encode encodes to jpg
func (e JpegEncoder) Encode(w io.Writer, img interface{}) error {
	m, ok := img.(*vips.ImageRef)
	if !ok {
		return errors.ErrInvalidType
	}

	buf, _, err := m.ExportJpeg(vips.NewJpegExportParams())
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

// Types returns the jpg suffixes.
func (e JpegEncoder) Types() []string {
	return []string{typeJpeg, typeJpg}
}

// MimeType returns the mimetype for jpg files.
func (e JpegEncoder) MimeType() string {
	return "image/jpeg"
}
