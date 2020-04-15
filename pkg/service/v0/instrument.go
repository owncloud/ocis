package svc

import (
	"github.com/owncloud/ocis-settings/pkg/metrics"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next Service, metrics *metrics.Metrics) Service {
	return Service{}
}

type instrument struct {
	next    Service
	metrics *metrics.Metrics
}
