// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/version"
	"github.com/owncloud/ocis/webdav/pkg/command"
	svcconfig "github.com/owncloud/ocis/webdav/pkg/config"
	"github.com/owncloud/ocis/webdav/pkg/flagset"
)

// WebDAVCommand is the entrypoint for the webdav command.
func WebDAVCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "webdav",
		Usage:    "Start webdav server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.WebDAV),
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.WebDAV),
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureWebDAV(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureWebDAV(cfg *config.Config) *svcconfig.Config {
	cfg.WebDAV.Log.Level = cfg.Log.Level
	cfg.WebDAV.Log.Pretty = cfg.Log.Pretty
	cfg.WebDAV.Log.Color = cfg.Log.Color
	cfg.WebDAV.Service.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.WebDAV.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.WebDAV.Tracing.Type = cfg.Tracing.Type
		cfg.WebDAV.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.WebDAV.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.WebDAV
}

func init() {
	register.AddCommand(WebDAVCommand)
}
