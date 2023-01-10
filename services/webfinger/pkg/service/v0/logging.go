package service

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/webfinger"
)

// NewLogging returns a service that logs messages.
func NewLogging(next Service, logger log.Logger) Service {
	return logging{
		next:   next,
		logger: logger,
	}
}

type logging struct {
	next   Service
	logger log.Logger
}

// Webfinger implements the Service interface.
func (l logging) Webfinger(ctx context.Context, resource, rel string) (webfinger.JSONResourceDescriptor, error) {
	l.logger.Debug().
		Str("resource", resource).
		Str("rel", rel).
		Msg("Webfinger")

	return l.next.Webfinger(ctx, resource, rel)
}
