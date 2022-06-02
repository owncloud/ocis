package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/storage-system/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageSystemCommand is the entrypoint for the StorageSystem command.
func StorageSystemCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.StorageSystem.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.StorageSystem.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.StorageSystem.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.StorageSystem),
	}
}

func init() {
	register.AddCommand(StorageSystemCommand)
}
