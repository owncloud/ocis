package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/webdav/pkg/command"
	"github.com/urfave/cli/v2"
)

// WebDAVCommand is the entrypoint for the webdav command.
func WebDAVCommand(cfg *config.Config) *cli.Command {

	return &cli.Command{
		Name:     "webdav",
		Usage:    "Start webdav server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.WebDAV),
	}
}

func init() {
	register.AddCommand(WebDAVCommand)
}
