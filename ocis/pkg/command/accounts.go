package command

import (
	"github.com/owncloud/ocis/accounts/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AccountsCommand is the entrypoint for the accounts command.
func AccountsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "accounts",
		Usage:    "Start accounts server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				return err
			}

			if cfg.Commons != nil {
				cfg.Accounts.Commons = cfg.Commons
			}

			return nil
		},
		Subcommands: command.GetCommands(cfg.Accounts),
	}
}

func init() {
	register.AddCommand(AccountsCommand)
}
