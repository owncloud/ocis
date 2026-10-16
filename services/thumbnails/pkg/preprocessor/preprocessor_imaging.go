//go:build !enable_vips

package preprocessor

import (
	"context"
	"io"

	"github.com/kovidgoyal/imaging"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ImageDecoder is a converter for the image file
type ImageDecoder struct{}

// Convert reads the image file and returns the thumbnail image
func (i ImageDecoder) Convert(ctx context.Context, r io.Reader) (interface{}, error) {
	span := trace.SpanFromContext(ctx)
	_, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameConvert,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.converter.provider", "imaging"),
			attribute.String("ocis.thumbnails.converter.type", "ImageDecoder"),
		),
	)
	defer newSpan.End()

	img, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return nil, errors.Wrap(err, `could not decode the image`)
	}
	return img, nil
}
