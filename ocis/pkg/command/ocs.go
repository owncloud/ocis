package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocs/pkg/command"
	"github.com/urfave/cli/v2"
)

// OCSCommand is the entrypoint for the ocs command.
func OCSCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "ocs",
		Usage:    "Start ocs server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				return err
			}

			if cfg.Commons != nil {
				cfg.OCS.Commons = cfg.Commons
			}

			return nil
		},
		Subcommands: command.GetCommands(cfg.OCS),
	}
}

func init() {
	register.AddCommand(OCSCommand)
}
