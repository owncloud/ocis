//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
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
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.WebDAV),
		},
		Before: func(ctx *cli.Context) error {
			if cfg.Commons != nil {
				cfg.WebDAV.Commons = cfg.Commons
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(cfg.WebDAV)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(WebDAVCommand)
}
