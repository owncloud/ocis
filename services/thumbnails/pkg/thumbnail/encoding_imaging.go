//go:build !enable_vips

package thumbnail

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"io"

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
			attribute.String("ocis.thumbnails.encoder.type", "PngEncoder"),
			attribute.String("ocis.thumbnails.encoder.mime", e.MimeType()),
		),
	)
	defer newSpan.End()

	m, ok := img.(image.Image)
	if !ok {
		return errors.ErrInvalidType
	}
	return png.Encode(w, m)
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
			attribute.String("ocis.thumbnails.encoder.type", "JpegEncoder"),
			attribute.String("ocis.thumbnails.encoder.mime", e.MimeType()),
		),
	)
	defer newSpan.End()

	m, ok := img.(image.Image)
	if !ok {
		return errors.ErrInvalidType
	}
	return jpeg.Encode(w, m, nil)
}

// Types returns the jpg suffixes.
func (e JpegEncoder) Types() []string {
	return []string{typeJpeg, typeJpg}
}

// MimeType returns the mimetype for jpg files.
func (e JpegEncoder) MimeType() string {
	return "image/jpeg"
}
