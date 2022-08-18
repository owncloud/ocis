package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/command"
	"github.com/urfave/cli/v2"
)

// ThumbnailsCommand is the entrypoint for the thumbnails command.
func ThumbnailsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Thumbnails.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Thumbnails.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.LogError(parser.ParseConfig(cfg, true))
			cfg.Thumbnails.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Thumbnails),
	}
}

func init() {
	register.AddCommand(ThumbnailsCommand)
}
