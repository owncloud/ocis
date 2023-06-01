package tracing

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/auth-service/pkg/config"
	"go.opentelemetry.io/otel/trace"

	pkgtrace "github.com/owncloud/ocis/v2/ocis-pkg/tracing"
)

var (
	// TraceProvider is the global trace provider for the service.
	TraceProvider = trace.NewNoopTracerProvider()
)

func Configure(cfg *config.Config, logger log.Logger) error {
	var err error
	if cfg.Tracing.Enabled {
		if TraceProvider, err = pkgtrace.GetTraceProvider(cfg.Tracing.Endpoint, cfg.Tracing.Collector, cfg.Service.Name, cfg.Tracing.Type); err != nil {
			return err
		}
	}

	return nil
}
