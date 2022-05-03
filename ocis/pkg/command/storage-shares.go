package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/storage-shares/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageSharesCommand is the entrypoint for the StorageShares command.
func StorageSharesCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.StorageShares.Service.Name,
		Usage:    subcommandDescription(cfg.StorageShares.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.StorageShares.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.StorageShares),
	}
}

func init() {
	register.AddCommand(StorageSharesCommand)
}
