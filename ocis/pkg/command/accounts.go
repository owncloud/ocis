package command

import (
	"github.com/owncloud/ocis/extensions/accounts/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AccountsCommand is the entrypoint for the accounts command.
func AccountsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        cfg.Accounts.Service.Name,
		Usage:       subcommandDescription(cfg.Accounts.Service.Name),
		Category:    "extensions",
		Subcommands: command.GetCommands(cfg.Accounts),
	}
}

func init() {
	register.AddCommand(AccountsCommand)
}
