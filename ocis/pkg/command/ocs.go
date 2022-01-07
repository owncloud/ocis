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
		Name:     cfg.OCS.Service.Name,
		Usage:    subcommandDescription(cfg.OCS.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.OCS),
	}
}

func init() {
	register.AddCommand(OCSCommand)
}
