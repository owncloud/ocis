package svc

import (
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

// NewLogging returns a service that logs messages.
func NewLoggingHandler(next Service, logger log.Logger) Service {
	return loggingHandler{
		next:   next,
		logger: logger,
	}
}

type loggingHandler struct {
	next   Service
	logger log.Logger
}

// ServeHTTP implements the Service interface.
func (l loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next.ServeHTTP(w, r)
}
