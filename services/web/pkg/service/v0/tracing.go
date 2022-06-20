package svc

import (
	"net/http"
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
	t.next.ServeHTTP(w, r)
}

// Config implements the Service interface.
func (t tracing) Config(w http.ResponseWriter, r *http.Request) {
	t.next.Config(w, r)
}
