package decorators

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	thumbnailssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/thumbnails/v0"
	"go.opentelemetry.io/otel/attribute"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next DecoratedService, tp trace.TracerProvider) DecoratedService {
	return tracing{
		Decorator: Decorator{next: next},
		tp:        tp,
	}
}

type tracing struct {
	Decorator
	tp trace.TracerProvider
}

// GetThumbnail implements the ThumbnailServiceHandler interface.
func (t tracing) GetThumbnail(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, rsp *thumbnailssvc.GetThumbnailResponse) error {
	var span trace.Span

	if t.tp != nil {
		tracer := t.tp.Tracer("thumbnails")
		spanOpts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
		}
		ctx, span = tracer.Start(ctx, "Thumbnails.GetThumbnail", spanOpts...)
		defer span.End()

		span.SetAttributes(
			attribute.KeyValue{Key: "filepath", Value: attribute.StringValue(req.GetFilepath())},
			attribute.KeyValue{Key: "thumbnail_type", Value: attribute.StringValue(req.GetThumbnailType().String())},
			attribute.KeyValue{Key: "width", Value: attribute.IntValue(int(req.GetWidth()))},
			attribute.KeyValue{Key: "height", Value: attribute.IntValue(int(req.GetHeight()))},
		)
	}

	return t.next.GetThumbnail(ctx, req, rsp)
}
