//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/version"
	"github.com/owncloud/ocis/store/pkg/command"
	svcconfig "github.com/owncloud/ocis/store/pkg/config"
	"github.com/owncloud/ocis/store/pkg/flagset"
	"github.com/urfave/cli/v2"
)

// StoreCommand is the entrypoint for the ocs command.
func StoreCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "store",
		Usage:    "Start a go-micro store",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Store),
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Store),
		},
		Before: func(ctx *cli.Context) error {
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureStore(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureStore(cfg *config.Config) *svcconfig.Config {
	cfg.Store.Log.Level = cfg.Log.Level
	cfg.Store.Log.Pretty = cfg.Log.Pretty
	cfg.Store.Log.Color = cfg.Log.Color
	cfg.Store.Service.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.Store.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Store.Tracing.Type = cfg.Tracing.Type
		cfg.Store.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Store.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.Store
}

func init() {
	register.AddCommand(StoreCommand)
}
