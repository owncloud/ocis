package svc

import (
	"net/http"

	"github.com/owncloud/ocis/ocs/pkg/metrics"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next Service, metrics *metrics.Metrics) Service {
	return instrument{
		next:    next,
		metrics: metrics,
	}
}

type instrument struct {
	next    Service
	metrics *metrics.Metrics
}

// ServeHTTP implements the Service interface.
func (i instrument) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.next.ServeHTTP(w, r)
}

// GetConfig implements the Service interface.
func (i instrument) GetConfig(w http.ResponseWriter, r *http.Request) {
	i.next.GetConfig(w, r)
}
