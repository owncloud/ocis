package svc

import (
	"github.com/owncloud/ocis/ocis-pkg/log"
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

// ListenForEvents implements service interface
func (l logging) ListenForEvents() {
	l.next.ListenForEvents()
}
