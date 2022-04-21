package command

import (
	"github.com/owncloud/ocis/extensions/storage-publiclink/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StoragePublicLinkCommand is the entrypoint for the reva-storage-oc command.
func StoragePublicLinkCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-public-link",
		Usage:    "start storage public link storage",
		Category: "extensions",
		//Flags:    flagset.StoragePublicLink(cfg.Storage),
		// Before: func(ctx *cli.Context) error {
		// 	return ParseStorageCommon(ctx, cfg)
		// },
		Action: func(c *cli.Context) error {
			origCmd := command.StoragePublicLink(cfg.StoragePublicLink)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StoragePublicLinkCommand)
}
