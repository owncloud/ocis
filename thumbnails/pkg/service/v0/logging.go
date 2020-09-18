package svc

import (
	"context"
	"time"

	"github.com/owncloud/ocis-pkg/v2/log"
	v0proto "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
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
	start := time.Now()
	err := l.next.GetThumbnail(ctx, req, rsp)

	logger := l.logger.With().
		Str("method", "Thumbnails.GetThumbnail").
		Dur("duration", time.Since(start)).
		Logger()

	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Failed to execute")
	} else {
		logger.Debug().
			Msg("")
	}
	return err
}
