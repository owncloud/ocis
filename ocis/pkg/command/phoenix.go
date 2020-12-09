package command

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/web/pkg/command"
	"github.com/owncloud/ocis/web/pkg/flagset"
)

// WebCommand is the entrypoint for the web command.
func WebCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "web",
		Usage:    "Start web server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Web),
		Before: func(c *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			cfg.Web.Web.Config.Apps = c.StringSlice("web-config-app")
			return nil
		},
		Action: func(c *cli.Context) error {
			webCommand := command.Server(configureWeb(cfg).Web)

			if err := webCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(webCommand.Action, c)
		},
	}
}

func configureWeb(cfg *config.Config) *config.Config {
	cfg.Web.Log.Level = cfg.Log.Level
	cfg.Web.Log.Pretty = cfg.Log.Pretty
	cfg.Web.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Web.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Web.Tracing.Type = cfg.Tracing.Type
		cfg.Web.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Web.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg
}

func init() {
	register.AddCommand(WebCommand)
}
