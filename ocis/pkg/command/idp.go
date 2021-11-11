package command

import (
	"github.com/owncloud/ocis/idp/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// IDPCommand is the entrypoint for the idp command.
func IDPCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "idp",
		Usage:    "Start idp server",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.IDP),
		},
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			if cfg.Commons != nil {
				cfg.IDP.Commons = cfg.Commons
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			idpCommand := command.Server(cfg.IDP)
			if err := idpCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(idpCommand.Action, c)
		},
	}
}

func init() {
	register.AddCommand(IDPCommand)
}
