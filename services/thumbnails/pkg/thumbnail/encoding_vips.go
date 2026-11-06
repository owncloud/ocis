//go:build enable_vips

package thumbnail

import (
	"context"
	"io"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// PngEncoder encodes to png
type PngEncoder struct{}

// Encode encodes to png format
func (e PngEncoder) Encode(ctx context.Context, w io.Writer, img interface{}) error {
	span := trace.SpanFromContext(ctx)
	_, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameEncode,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.encoder.type", "PngEncoderVips"),
			attribute.String("ocis.thumbnails.encoder.mime", e.MimeType()),
		),
	)
	defer newSpan.End()

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
func (e JpegEncoder) Encode(ctx context.Context, w io.Writer, img interface{}) error {
	span := trace.SpanFromContext(ctx)
	_, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameEncode,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.encoder.type", "JpegEncoderVips"),
			attribute.String("ocis.thumbnails.encoder.mime", e.MimeType()),
		),
	)
	defer newSpan.End()

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
