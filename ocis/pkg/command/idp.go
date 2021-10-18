package command

import (
	"github.com/owncloud/ocis/idp/pkg/command"
	svcconfig "github.com/owncloud/ocis/idp/pkg/config"
	"github.com/owncloud/ocis/idp/pkg/flagset"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// IDPCommand is the entrypoint for the idp command.
func IDPCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "idp",
		Usage:    "Start idp server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.IDP),
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.IDP),
		},
		Before: func(ctx *cli.Context) error {
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			idpCommand := command.Server(configureIDP(cfg))

			if err := idpCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(idpCommand.Action, c)
		},
	}
}

func configureIDP(cfg *config.Config) *svcconfig.Config {
	cfg.IDP.Log.Level = cfg.Log.Level
	cfg.IDP.Log.Pretty = cfg.Log.Pretty
	cfg.IDP.Log.Color = cfg.Log.Color
	cfg.IDP.HTTP.TLS = false
	cfg.IDP.Service.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.IDP.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.IDP.Tracing.Type = cfg.Tracing.Type
		cfg.IDP.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.IDP.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.IDP
}

func init() {
	register.AddCommand(IDPCommand)
}
