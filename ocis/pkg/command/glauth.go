package command

import (
	"github.com/owncloud/ocis/glauth/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GLAuthCommand is the entrypoint for the glauth command.
func GLAuthCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.GLAuth.Service.Name,
		Usage:    subcommandDescription(cfg.GLAuth.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.GLAuth),
	}
}

func init() {
	register.AddCommand(GLAuthCommand)
}
