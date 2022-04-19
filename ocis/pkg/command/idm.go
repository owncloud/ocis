package command

import (
	"github.com/owncloud/ocis/extensions/idm/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// IDMCommand is the entrypoint for the idm server command.
func IDMCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "idm",
		Usage:    "idm extension commands",
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.IDM),
	}
}

func init() {
	register.AddCommand(IDMCommand)
}
