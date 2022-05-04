package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/accounts/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AccountsCommand is the entrypoint for the accounts command.
func AccountsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Accounts.Service.Name,
		Usage:    subcommandDescription(cfg.Accounts.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Accounts.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Accounts),
	}
}

func init() {
	register.AddCommand(AccountsCommand)
}
