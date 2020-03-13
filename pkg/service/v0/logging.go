package svc

import (
	"context"

	"github.com/owncloud/ocis-pkg/v2/log"
	v0proto "github.com/owncloud/ocis-thumbnails/pkg/proto/v0"
)

// NewLogging returns a service that logs messages.
func NewLogging(next v0proto.ThumbnailServiceHandler, logger log.Logger) v0proto.ThumbnailServiceHandler {
	return logging{
		next:   next,
		logger: logger,
	}
}

type logging struct {
	next   v0proto.ThumbnailServiceHandler
	logger log.Logger
}

// GetThumbnail implements the ThumbnailServiceHandler interface.
func (l logging) GetThumbnail(ctx context.Context, req *v0proto.GetRequest, rsp *v0proto.GetResponse) error {
	return l.next.GetThumbnail(ctx, req, rsp)
}
