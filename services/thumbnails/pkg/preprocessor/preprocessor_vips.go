//go:build enable_vips

package preprocessor

import (
	"context"
	"io"

	"github.com/davidbyttow/govips/v2/vips"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	vips.LoggingSettings(nil, vips.LogLevelError)
}

type ImageDecoder struct{}

func (v ImageDecoder) Convert(ctx context.Context, r io.Reader) (interface{}, error) {
	span := trace.SpanFromContext(ctx)
	_, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameConvert,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.converter.provider", "vips"),
			attribute.String("ocis.thumbnails.converter.type", "ImageDecoder"),
		),
	)
	defer newSpan.End()

	img, err := vips.NewImageFromReader(r)
	return img, err
}
