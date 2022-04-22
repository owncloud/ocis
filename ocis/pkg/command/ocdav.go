package command

import (
	"github.com/owncloud/ocis/extensions/ocdav/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// OCDavCommand is the entrypoint for the ocdav command.
func OCDavCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "ocdav",
		Usage:    "start ocdav",
		Category: "extensions",
		// Before: func(ctx *cli.Context) error {
		// 	return ParseStorageCommon(ctx, cfg)
		// },
		Action: func(c *cli.Context) error {
			origCmd := command.OCDav(cfg.OCDav)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(OCDavCommand)
}
