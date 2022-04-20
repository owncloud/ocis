package command

import (
	"github.com/owncloud/ocis/extensions/sharing/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageSharingCommand is the entrypoint for the reva-sharing command.
func StorageSharingCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-sharing",
		Usage:    "start storage sharing service",
		Category: "extensions",
		//Flags:    flagset.SharingWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
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
