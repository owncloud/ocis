package decorators

import (
	"context"
	"net/http"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	thumbnailssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/thumbnails/v0"
	merrors "go-micro.dev/v4/errors"
)

// NewLogging returns a service that logs messages.
func NewLogging(next DecoratedService, logger log.Logger) DecoratedService {
	return logging{
		Decorator: Decorator{next: next},
		logger:    logger,
	}
}

type logging struct {
	Decorator
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
		fromError := merrors.FromError(err)
		switch fromError.GetCode() {
		case http.StatusNotFound:
			logger.Debug().
				Str("error_detail", fromError.GetDetail()).
				Msg("no thumbnail found")
		default:
			logger.Warn().
				Err(err).
				Msg("Failed to execute")
		}
	} else {
		logger.Debug().
			Msg("")
	}
	return err
}
