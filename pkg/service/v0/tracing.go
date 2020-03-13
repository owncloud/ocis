package svc

import (
	"context"

	v0proto "github.com/owncloud/ocis-thumbnails/pkg/proto/v0"
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
func (t tracing) GetThumbnail(ctx context.Context, req *v0proto.GetRequest, rsp *v0proto.GetResponse) error {
	return t.next.GetThumbnail(ctx, req, rsp)
}
