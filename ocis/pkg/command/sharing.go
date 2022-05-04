package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/sharing/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// SharingCommand is the entrypoint for the Sharing command.
func SharingCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Sharing.Service.Name,
		Usage:    subcommandDescription(cfg.Sharing.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Sharing.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Sharing),
	}
}

func init() {
	register.AddCommand(SharingCommand)
}
