package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/storage-shares/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageSharesCommand is the entrypoint for the StorageShares command.
func StorageSharesCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.StorageShares.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.StorageShares.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
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
