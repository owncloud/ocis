package svc

import (
	"context"

	v0proto "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"go.opencensus.io/trace"
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
	ctx, span := trace.StartSpan(ctx, "Thumbnails.GetThumbnail")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("filepath", req.Filepath),
		trace.StringAttribute("thumbnail_type", req.ThumbnailType.String()),
		trace.Int64Attribute("width", int64(req.Width)),
		trace.Int64Attribute("height", int64(req.Height)),
	}, "Execute Thumbnails.GetThumbnail handler")

	return t.next.GetThumbnail(ctx, req, rsp)
}
