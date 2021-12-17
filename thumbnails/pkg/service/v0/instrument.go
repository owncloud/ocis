package svc

import (
	"context"

	thumbnailssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/thumbnails/v1"
	"github.com/owncloud/ocis/thumbnails/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next thumbnailssvc.ThumbnailServiceHandler, metrics *metrics.Metrics) thumbnailssvc.ThumbnailServiceHandler {
	return instrument{
		next:    next,
		metrics: metrics,
	}
}

type instrument struct {
	next    thumbnailssvc.ThumbnailServiceHandler
	metrics *metrics.Metrics
}

// GetThumbnail implements the ThumbnailServiceHandler interface.
func (i instrument) GetThumbnail(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, rsp *thumbnailssvc.GetThumbnailResponse) error {
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
