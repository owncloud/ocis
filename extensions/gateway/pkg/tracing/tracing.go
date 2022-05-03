package tracing

import (
	"github.com/owncloud/ocis/extensions/gateway/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

var (
	// TraceProvider is the global trace provider for the proxy service.
	TraceProvider = trace.NewNoopTracerProvider()
)

func Configure(cfg *config.Config, logger log.Logger) error {
	tracing.Configure(cfg.Tracing.Enabled, cfg.Tracing.Type, logger)
	return nil
}
