package tracing

import (
	"github.com/owncloud/ocis/v2/extensions/proxy/pkg/config"
	pkgtrace "github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

var (
	// TraceProvider is the global trace provider for the proxy service.
	TraceProvider = trace.NewNoopTracerProvider()
)

func Configure(cfg *config.Config) error {
	var err error
	if cfg.Tracing.Enabled {
		if TraceProvider, err = pkgtrace.GetTraceProvider(cfg.Tracing.Endpoint, cfg.Tracing.Collector, cfg.Service.Name, cfg.Tracing.Type); err != nil {
			return err
		}
	}

	return nil
}
