package svc

import (
	"context"

	"github.com/owncloud/ocis-thumbnails/pkg/metrics"
	v0proto "github.com/owncloud/ocis-thumbnails/pkg/proto/v0"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next v0proto.ThumbnailServiceHandler, metrics *metrics.Metrics) v0proto.ThumbnailServiceHandler {
	return instrument{
		next:    next,
		metrics: metrics,
	}
}

type instrument struct {
	next    v0proto.ThumbnailServiceHandler
	metrics *metrics.Metrics
}

// GetThumbnail implements the ThumbnailServiceHandler interface.
func (i instrument) GetThumbnail(ctx context.Context, req *v0proto.GetRequest, rsp *v0proto.GetResponse) error {
	return i.next.GetThumbnail(ctx, req, rsp)
}
