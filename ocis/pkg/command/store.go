package command

import (
	"github.com/owncloud/ocis/extensions/store/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StoreCommand is the entrypoint for the ocs command.
func StoreCommand(cfg *config.Config) *cli.Command {

	return &cli.Command{
		Name:        cfg.Store.Service.Name,
		Usage:       subcommandDescription(cfg.Store.Service.Name),
		Category:    "extensions",
		Subcommands: command.GetCommands(cfg.Store),
	}
}

func init() {
	register.AddCommand(StoreCommand)
}
