// +build !simple

package command

import (
	svcconfig "github.com/owncloud/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis/pkg/config"
)

func configurePhoenix(cfg *config.Config) *svcconfig.Config {
	cfg.Phoenix.Log.Level = cfg.Log.Level
	cfg.Phoenix.Log.Pretty = cfg.Log.Pretty
	cfg.Phoenix.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Phoenix.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Phoenix.Tracing.Type = cfg.Tracing.Type
		cfg.Phoenix.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Phoenix.Tracing.Collector = cfg.Tracing.Collector
		cfg.Phoenix.Tracing.Service = cfg.Tracing.Service
	}

	// disable ocis-hello extension
	cfg.Phoenix.Phoenix.Config.ExternalApps = []svcconfig.ExternalApp{}

	return cfg.Phoenix
}
