//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/accounts/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AccountsCommand is the entrypoint for the accounts command.
func AccountsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "accounts",
		Usage:    "Start accounts server",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.ListAccounts(cfg.Accounts),
			command.AddAccount(cfg.Accounts),
			command.UpdateAccount(cfg.Accounts),
			command.RemoveAccount(cfg.Accounts),
			command.InspectAccount(cfg.Accounts),
			command.PrintVersion(cfg.Accounts),
		},
		Before: func(ctx *cli.Context) error {
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			cfg.Accounts.Log = cfg.Log
			origCmd := command.Server(cfg.Accounts)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(AccountsCommand)
}
