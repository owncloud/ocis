//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/thumbnails/pkg/command"
	"github.com/urfave/cli/v2"
)

// ThumbnailsCommand is the entrypoint for the thumbnails command.
func ThumbnailsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "thumbnails",
		Usage:    "Start thumbnails server",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Thumbnails),
		},
		Before: func(ctx *cli.Context) error {
			if cfg.Commons != nil {
				cfg.Thumbnails.Commons = cfg.Commons
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(cfg.Thumbnails)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(ThumbnailsCommand)
}
