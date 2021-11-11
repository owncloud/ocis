//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/urfave/cli/v2"
)

// StoragePublicLinkCommand is the entrypoint for the reva-storage-oc command.
func StoragePublicLinkCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-public-link",
		Usage:    "Start storage public link storage",
		Category: "Extensions",
		Flags:    flagset.StoragePublicLink(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.StoragePublicLink(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StoragePublicLinkCommand)
}
