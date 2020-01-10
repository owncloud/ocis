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
	cfg.Phoenix.Tracing.Enabled = false
	cfg.Phoenix.HTTP.Addr = "localhost:9100"
	cfg.Phoenix.HTTP.Root = "/"
	cfg.Phoenix.Phoenix.Config.Apps = []string{
		"draw-io",
		"files",
		"markdown-viewer",
		"media-viewer",
		"pdf-viewer",
	}

	return cfg.Phoenix
}
