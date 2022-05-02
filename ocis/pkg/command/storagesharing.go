package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/sharing/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageSharingCommand is the entrypoint for the reva-sharing command.
func StorageSharingCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-sharing",
		Usage:    "start storage sharing service",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Sharing.Commons = cfg.Commons
			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Sharing(cfg.Sharing)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageSharingCommand)
}
