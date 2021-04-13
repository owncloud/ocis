package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/web/pkg/command"
	svcconfig "github.com/owncloud/ocis/web/pkg/config"
	"github.com/owncloud/ocis/web/pkg/flagset"
)

// WebCommand is the entrypoint for the web command.
func WebCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "web",
		Usage:    "Start web server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Web),
		Before: func(ctx *cli.Context) error {
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureWeb(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureWeb(cfg *config.Config) *svcconfig.Config {
	cfg.Web.Log.Level = cfg.Log.Level
	cfg.Web.Log.Pretty = cfg.Log.Pretty
	cfg.Web.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Web.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Web.Tracing.Type = cfg.Tracing.Type
		cfg.Web.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Web.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.Web
}

func init() {
	register.AddCommand(WebCommand)
}
