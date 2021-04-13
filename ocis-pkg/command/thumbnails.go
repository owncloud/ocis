// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/version"
	"github.com/owncloud/ocis/thumbnails/pkg/command"
	"github.com/owncloud/ocis/thumbnails/pkg/flagset"

	svcconfig "github.com/owncloud/ocis/thumbnails/pkg/config"
)

// ThumbnailsCommand is the entrypoint for the thumbnails command.
func ThumbnailsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "thumbnails",
		Usage:    "Start thumbnails server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Thumbnails),
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Thumbnails),
		},
		Before: func(ctx *cli.Context) error {
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureThumbnails(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureThumbnails(cfg *config.Config) *svcconfig.Config {
	cfg.Thumbnails.Log.Level = cfg.Log.Level
	cfg.Thumbnails.Log.Pretty = cfg.Log.Pretty
	cfg.Thumbnails.Log.Color = cfg.Log.Color
	cfg.Thumbnails.Server.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.Thumbnails.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Thumbnails.Tracing.Type = cfg.Tracing.Type
		cfg.Thumbnails.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Thumbnails.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.Thumbnails
}

func init() {
	register.AddCommand(ThumbnailsCommand)
}
