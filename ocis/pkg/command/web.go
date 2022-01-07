package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/web/pkg/command"
	"github.com/urfave/cli/v2"
)

// WebCommand is the entrypoint for the web command.
func WebCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Web.Service.Name,
		Usage:    subcommandDescription(cfg.Web.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.Web),
	}
}

func init() {
	register.AddCommand(WebCommand)
}
