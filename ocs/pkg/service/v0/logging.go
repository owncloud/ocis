package svc

import (
	"net/http"

	"github.com/owncloud/ocis-pkg/v2/log"
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

// ServeHTTP implements the Service interface.
func (l logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next.ServeHTTP(w, r)
}

// GetConfig implements the Service interface.
func (l logging) GetConfig(w http.ResponseWriter, r *http.Request) {
	l.next.GetConfig(w, r)
}
