package svc

import (
	"context"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/log"
	thumbnailssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/thumbnails/v0"
)

// NewLogging returns a service that logs messages.
func NewLogging(next thumbnailssvc.ThumbnailServiceHandler, logger log.Logger) thumbnailssvc.ThumbnailServiceHandler {
	return logging{
		next:   next,
		logger: logger,
	}
}

type logging struct {
	next   thumbnailssvc.ThumbnailServiceHandler
	logger log.Logger
}

// GetThumbnail implements the ThumbnailServiceHandler interface.
func (l logging) GetThumbnail(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, rsp *thumbnailssvc.GetThumbnailResponse) error {
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
