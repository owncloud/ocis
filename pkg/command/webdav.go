// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-webdav/pkg/command"
	svcconfig "github.com/owncloud/ocis-webdav/pkg/config"
	"github.com/owncloud/ocis-webdav/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// WebDAVCommand is the entrypoint for the webdav command.
func WebDAVCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "webdav",
		Usage:    "Start webdav server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.WebDAV),
		Action: func(c *cli.Context) error {
			scfg := configureWebDAV(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureWebDAV(cfg *config.Config) *svcconfig.Config {
	cfg.WebDAV.Log.Level = cfg.Log.Level
	cfg.WebDAV.Log.Pretty = cfg.Log.Pretty
	cfg.WebDAV.Log.Color = cfg.Log.Color
	cfg.WebDAV.Tracing.Enabled = false
	cfg.WebDAV.HTTP.Addr = "localhost:9115"
	cfg.WebDAV.HTTP.Root = "/"

	return cfg.WebDAV
}

func init() {
	register.AddCommand(WebDAVCommand)
}
