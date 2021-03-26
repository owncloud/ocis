package tracing

import (
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// Configure for Reva serves only as informational / instructive log messages. Tracing config will be delegated directly
// to Reva services.
func Configure(cfg *config.Config, logger log.Logger) {
	if cfg.Tracing.Enabled {
		switch t := cfg.Tracing.Type; t {
		case "agent":
			logger.Error().
				Str("type", t).
				Msg("Reva only supports the jaeger tracing backend")

		case "jaeger":
			logger.Info().
				Str("type", t).
				Msg("configuring storage to use the jaeger tracing backend")

		case "zipkin":
			logger.Error().
				Str("type", t).
				Msg("Reva only supports the jaeger tracing backend")

		default:
			logger.Warn().
				Str("type", t).
				Msg("Unknown tracing backend")
		}

	} else {
		logger.Debug().
			Msg("Tracing is not enabled")
	}
}
