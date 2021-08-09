package tracing

import (
	"github.com/owncloud/ocis/accounts/pkg/config"
	pkgtrace "github.com/owncloud/ocis/ocis-pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

var (
	// TraceProvider is the global trace provider for the proxy service.
	TraceProvider = trace.NewNoopTracerProvider()
)

func Configure(cfg *config.Config) error {
	var err error
	if TraceProvider, err = pkgtrace.GetTraceProvider(cfg.Tracing.Collector, cfg.Tracing.Type, "accounts"); err != nil {
		return err
	}

	return nil
}
