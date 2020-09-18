package command

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-phoenix/pkg/command"
	"github.com/owncloud/ocis/ocis-phoenix/pkg/flagset"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
)

// PhoenixCommand is the entrypoint for the phoenix command.
func PhoenixCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "phoenix",
		Usage:    "Start phoenix server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Phoenix),
		Before: func(c *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			cfg.Phoenix.Phoenix.Config.Apps = c.StringSlice("web-config-app")
			return nil
		},
		Action: func(c *cli.Context) error {
			phoenixCommand := command.Server(configurePhoenix(cfg).Phoenix)

			if err := phoenixCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(phoenixCommand.Action, c)
		},
	}
}

func configurePhoenix(cfg *config.Config) *config.Config {
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

	return cfg
}

func init() {
	register.AddCommand(PhoenixCommand)
}
