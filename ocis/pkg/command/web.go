package command

import (
	"github.com/owncloud/ocis/extensions/web/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// WebCommand is the entrypoint for the web command.
func WebCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        cfg.Web.Service.Name,
		Usage:       subcommandDescription(cfg.Web.Service.Name),
		Category:    "extensions",
		Subcommands: command.GetCommands(cfg.Web),
	}
}

func init() {
	register.AddCommand(WebCommand)
}
