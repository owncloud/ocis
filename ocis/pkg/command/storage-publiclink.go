package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/storage-publiclink/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StoragePublicLinkCommand is the entrypoint for the StoragePublicLink command.
func StoragePublicLinkCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.StoragePublicLink.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.StoragePublicLink.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.StoragePublicLink.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.StoragePublicLink),
	}
}

func init() {
	register.AddCommand(StoragePublicLinkCommand)
}
