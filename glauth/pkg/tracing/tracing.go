package tracing

import (
	"github.com/owncloud/ocis/glauth/pkg/config"
	pkgtrace "github.com/owncloud/ocis/ocis-pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

var (
	// TraceProvider is the global trace provider for the glauth service.
	TraceProvider = trace.NewNoopTracerProvider()
)

func Configure(cfg *config.Config) error {
	var err error
	if TraceProvider, err = pkgtrace.GetTraceProvider(cfg.Tracing.Endpoint, cfg.Tracing.Collector, "glauth", cfg.Tracing.Type); err != nil {
		return err
	}

	return nil
}
