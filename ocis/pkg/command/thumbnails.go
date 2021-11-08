//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/thumbnails/pkg/command"
	"github.com/urfave/cli/v2"
)

// ThumbnailsCommand is the entrypoint for the thumbnails command.
func ThumbnailsCommand(cfg *config.Config) *cli.Command {
	var globalLog shared.Log

	return &cli.Command{
		Name:     "thumbnails",
		Usage:    "Start thumbnails server",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Thumbnails),
		},
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			globalLog = cfg.Log

			return nil
		},
		Action: func(c *cli.Context) error {
			// if thumbnails logging is empty in ocis.yaml
			if (cfg.Thumbnails.Log == shared.Log{}) && (globalLog != shared.Log{}) {
				// we can safely inherit the global logging values.
				cfg.Thumbnails.Log = globalLog
			}
			origCmd := command.Server(cfg.Thumbnails)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(ThumbnailsCommand)
}
