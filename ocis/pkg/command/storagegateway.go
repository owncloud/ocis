//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageGatewayCommand is the entrypoint for the reva-gateway command.
func StorageGatewayCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-gateway",
		Usage:    "Start storage gateway",
		Category: "Extensions",
		//Flags:    flagset.GatewayWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Gateway(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageGatewayCommand)
}
