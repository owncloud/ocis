package tracing

import (
	"github.com/owncloud/ocis/extensions/storage/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Configure for Reva serves only as informational / instructive log messages. Tracing config will be delegated directly
// to Reva services.
func Configure(cfg *config.Config, logger log.Logger) {
	if cfg.Tracing.Enabled {
		switch cfg.Tracing.Type {
		case "agent":
			logger.Error().
				Str("type", cfg.Tracing.Type).
				Msg("Reva only supports the jaeger tracing backend")

		case "jaeger":
			logger.Info().
				Str("type", cfg.Tracing.Type).
				Msg("configuring storage to use the jaeger tracing backend")

		case "zipkin":
			logger.Error().
				Str("type", cfg.Tracing.Type).
				Msg("Reva only supports the jaeger tracing backend")

		default:
			logger.Warn().
				Str("type", cfg.Tracing.Type).
				Msg("Unknown tracing backend")
		}

	} else {
		logger.Debug().
			Msg("Tracing is not enabled")
	}
}
