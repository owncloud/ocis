package svc

import (
	"net/http"

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

// ServeHTTP implements the Service interface.
func (l logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next.ServeHTTP(w, r)
}

// Dummy implements the Service interface.
func (l logging) Dummy(w http.ResponseWriter, r *http.Request) {
	l.next.Dummy(w, r)
}
