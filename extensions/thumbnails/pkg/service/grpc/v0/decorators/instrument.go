package decorators

import (
	"context"

	"github.com/owncloud/ocis/extensions/thumbnails/pkg/metrics"
	thumbnailssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/thumbnails/v0"
	"github.com/prometheus/client_golang/prometheus"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next DecoratedService, metrics *metrics.Metrics) DecoratedService {
	return instrument{
		Decorator: Decorator{next: next},
		metrics:   metrics,
	}
}

type instrument struct {
	Decorator
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
