package tracing

import (
	"github.com/owncloud/ocis/glauth/pkg/config"
	pkgtrace "github.com/owncloud/ocis/ocis-pkg/tracing"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	// TraceProvider is the global trace provider for the glauth service.
	TraceProvider = sdktrace.NewTracerProvider()
)

func Configure(cfg *config.Config) error {
	var err error
	if TraceProvider, err = pkgtrace.GetTraceProvider(cfg.Tracing.Collector, cfg.Tracing.Type, "glauth"); err != nil {
		return err
	}

	return nil
}
