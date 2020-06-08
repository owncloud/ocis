// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-accounts/pkg/command"
	svcconfig "github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// AccountsCommand is the entrypoint for the accounts command.
func AccountsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "accounts",
		Usage:    "Start accounts server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Accounts),
		Action: func(c *cli.Context) error {
			accountsCommand := command.Server(configureAccounts(cfg))
			if err := accountsCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(accountsCommand.Action, c)
		},
	}
}

func configureAccounts(cfg *config.Config) *svcconfig.Config {
	cfg.Accounts.Log.Level = cfg.Log.Level
	cfg.Accounts.Log.Pretty = cfg.Log.Pretty
	cfg.Accounts.Log.Color = cfg.Log.Color

	return cfg.Accounts
}

func init() {
	register.AddCommand(AccountsCommand)
}
