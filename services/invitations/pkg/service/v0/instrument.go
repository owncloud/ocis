package service

import (
	"context"

	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
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

// Invite implements the Service interface.
func (i instrument) Invite(ctx context.Context, invitation *invitations.Invitation) (*invitations.Invitation, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000

		i.metrics.Latency.WithLabelValues().Observe(us)
		i.metrics.Duration.WithLabelValues().Observe(v)
	}))

	defer timer.ObserveDuration()

	i.metrics.Counter.WithLabelValues().Inc()

	return i.next.Invite(ctx, invitation)
}
