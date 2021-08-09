package tracing

import (
	pkgtrace "github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/owncloud/ocis/proxy/pkg/config"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	// TraceProvider is the global trace provider for the proxy service.
	TraceProvider = sdktrace.NewTracerProvider()
)

func Configure(cfg *config.Config) error {
	var err error
	if TraceProvider, err = pkgtrace.GetTraceProvider(cfg.Tracing.Collector, cfg.Tracing.Type, "proxy"); err != nil {
		return err
	}

	return nil
}
