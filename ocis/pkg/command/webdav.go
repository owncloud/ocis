//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/webdav/pkg/command"
	"github.com/urfave/cli/v2"
)

// WebDAVCommand is the entrypoint for the webdav command.
func WebDAVCommand(cfg *config.Config) *cli.Command {
	var globalLog shared.Log

	return &cli.Command{
		Name:     "webdav",
		Usage:    "Start webdav server",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.WebDAV),
		},
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			globalLog = cfg.Log

			return nil
		},
		Action: func(c *cli.Context) error {
			// if webdav logging is empty in ocis.yaml
			if (cfg.WebDAV.Log == shared.Log{}) && (globalLog != shared.Log{}) {
				// we can safely inherit the global logging values.
				cfg.WebDAV.Log = globalLog
			}
			origCmd := command.Server(cfg.WebDAV)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(WebDAVCommand)
}
