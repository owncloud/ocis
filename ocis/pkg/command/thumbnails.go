package command

import (
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// ThumbnailsCommand is the entrypoint for the thumbnails command.
func ThumbnailsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        cfg.Thumbnails.Service.Name,
		Usage:       subcommandDescription(cfg.Thumbnails.Service.Name),
		Category:    "extensions",
		Subcommands: command.GetCommands(cfg.Thumbnails),
	}
}

func init() {
	register.AddCommand(ThumbnailsCommand)
}
