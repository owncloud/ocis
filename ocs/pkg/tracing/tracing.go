package tracing

import (
	pkgtrace "github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/owncloud/ocis/ocs/pkg/config"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	// TraceProvider is the global trace provider for the ocs service.
	TraceProvider = sdktrace.NewTracerProvider()
)

func Configure(cfg *config.Config) error {
	var err error
	if TraceProvider, err = pkgtrace.GetTraceProvider(cfg.Tracing.Collector, cfg.Tracing.Type, "ocs"); err != nil {
		return err
	}

	return nil
}
