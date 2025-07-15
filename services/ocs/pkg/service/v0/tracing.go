package svc

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next Service) Service {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next Service
}

// ServeHTTP implements the Service interface.
func (t tracing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.TraceContext(t.next).ServeHTTP(w, r)
}
