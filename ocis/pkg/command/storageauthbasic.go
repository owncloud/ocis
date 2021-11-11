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

// StorageAuthBasicCommand is the entrypoint for the reva-auth-basic command.
func StorageAuthBasicCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-auth-basic",
		Usage:    "Start storage auth-basic service",
		Category: "Extensions",
		Flags:    flagset.AuthBasicWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.AuthBasic(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAuthBasicCommand)
}
