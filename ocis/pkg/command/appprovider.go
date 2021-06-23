// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	svcconfig "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
)

// AppProviderCommand is the entrypoint for the reva-gateway command.
func AppProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "app-provider",
		Usage:    "Start appprovider for providing apps",
		Category: "Extensions",
		Flags:    flagset.AppProviderWithConfig(cfg.Storage),
		Action: func(c *cli.Context) error {
			origCmd := command.AppProvider(configureAppProvider(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureAppProvider(cfg *config.Config) *svcconfig.Config {
	cfg.Storage.Log.Level = cfg.Log.Level
	cfg.Storage.Log.Pretty = cfg.Log.Pretty
	cfg.Storage.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Storage.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Storage.Tracing.Type = cfg.Tracing.Type
		cfg.Storage.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Storage.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.Storage
}

func init() {
	register.AddCommand(AppProviderCommand)
}
