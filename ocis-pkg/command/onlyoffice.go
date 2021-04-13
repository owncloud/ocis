package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/onlyoffice/pkg/command"
	svcconfig "github.com/owncloud/ocis/onlyoffice/pkg/config"
	"github.com/owncloud/ocis/onlyoffice/pkg/flagset"
)

// OnlyofficeCommand is the entrypoint for the onlyoffice command.
func OnlyofficeCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "onlyoffice",
		Usage:    "Start onlyoffice server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Onlyoffice),
		Before: func(ctx *cli.Context) error {
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureOnlyoffice(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureOnlyoffice(cfg *config.Config) *svcconfig.Config {
	cfg.Onlyoffice.Log.Level = cfg.Log.Level
	cfg.Onlyoffice.Log.Pretty = cfg.Log.Pretty
	cfg.Onlyoffice.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Onlyoffice.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Onlyoffice.Tracing.Type = cfg.Tracing.Type
		cfg.Onlyoffice.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Onlyoffice.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.Onlyoffice
}

func init() {
	register.AddCommand(OnlyofficeCommand)
}
