package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/storage-users/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageUsersCommand is the entrypoint for the StorageUsers command.
func StorageUsersCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.StorageUsers.Service.Name,
		Usage:    subcommandDescription(cfg.StorageUsers.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.StorageUsers.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.StorageUsers),
	}
}

func init() {
	register.AddCommand(StorageUsersCommand)
}
