package svc

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	v0proto "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	thumbnailsTracing "github.com/owncloud/ocis/thumbnails/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next v0proto.ThumbnailServiceHandler) v0proto.ThumbnailServiceHandler {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next v0proto.ThumbnailServiceHandler
}

// GetThumbnail implements the ThumbnailServiceHandler interface.
func (t tracing) GetThumbnail(ctx context.Context, req *v0proto.GetThumbnailRequest, rsp *v0proto.GetThumbnailResponse) error {
	var span trace.Span

	if thumbnailsTracing.TraceProvider != nil {
		tracer := thumbnailsTracing.TraceProvider.Tracer("thumbnails")
		ctx, span = tracer.Start(ctx, "Thumbnails.GetThumbnail")
		defer span.End()

		span.SetAttributes(
			attribute.KeyValue{Key: "filepath", Value: attribute.StringValue(req.Filepath)},
			attribute.KeyValue{Key: "thumbnail_type", Value: attribute.StringValue(req.ThumbnailType.String())},
			attribute.KeyValue{Key: "width", Value: attribute.IntValue(int(req.Width))},
			attribute.KeyValue{Key: "height", Value: attribute.IntValue(int(req.Height))},
		)
	}

	return t.next.GetThumbnail(ctx, req, rsp)
}
