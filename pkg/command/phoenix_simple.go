// +build simple

package command

import (
	svcconfig "github.com/owncloud/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis/pkg/config"
)

func configurePhoenix(cfg *config.Config) *svcconfig.Config {
	cfg.Phoenix.Log.Level = cfg.Log.Level
	cfg.Phoenix.Log.Pretty = cfg.Log.Pretty
	cfg.Phoenix.Log.Color = cfg.Log.Color
	cfg.Phoenix.Tracing.Enabled = false
	cfg.Phoenix.HTTP.Addr = "localhost:9100"
	cfg.Phoenix.HTTP.Root = "/"
	// disable built in apps
	cfg.Phoenix.Phoenix.Config.Apps = []string{}
	// enable ocis-hello extension
	cfg.Phoenix.Phoenix.Config.ExternalApps = []svcconfig.ExternalApp{
		svcconfig.ExternalApp{
			Name: "hello",
			Path: "http://localhost:9105/hello.js",
			Config: map[string]interface{}{
				"url": "http://localhost:9105",
			},
		},
	}

	return cfg.Phoenix
}
