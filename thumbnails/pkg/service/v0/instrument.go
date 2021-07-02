package svc

import (
	"context"

	"github.com/owncloud/ocis/thumbnails/pkg/metrics"
	v0proto "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/prometheus/client_golang/prometheus"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next v0proto.ThumbnailServiceHandler, metrics *metrics.Metrics) v0proto.ThumbnailServiceHandler {
	return instrument{
		next:    next,
		metrics: metrics,
	}
}

type instrument struct {
	next    v0proto.ThumbnailServiceHandler
	metrics *metrics.Metrics
}

// GetThumbnail implements the ThumbnailServiceHandler interface.
func (i instrument) GetThumbnail(ctx context.Context, req *v0proto.GetThumbnailRequest, rsp *v0proto.GetThumbnailResponse) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000_000
		i.metrics.Latency.WithLabelValues().Observe(us)
		i.metrics.Duration.WithLabelValues().Observe(v)
	}))
	defer timer.ObserveDuration()

	err := i.next.GetThumbnail(ctx, req, rsp)

	if err != nil {
		i.metrics.Counter.WithLabelValues().Inc()
	}
	return err
}
