// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/version"
	"github.com/owncloud/ocis/proxy/pkg/command"
	svcconfig "github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/flagset"
)

// ProxyCommand is the entry point for the proxy command.
func ProxyCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "proxy",
		Usage:    "Start proxy server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Proxy),
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Proxy),
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureProxy(cfg))

			if err := origCmd.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(origCmd.Action, c)
		},
	}
}

func configureProxy(cfg *config.Config) *svcconfig.Config {
	cfg.Proxy.Log.Level = cfg.Log.Level
	cfg.Proxy.Log.Pretty = cfg.Log.Pretty
	cfg.Proxy.Log.Color = cfg.Log.Color
	cfg.Proxy.Service.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.Proxy.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Proxy.Tracing.Type = cfg.Tracing.Type
		cfg.Proxy.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Proxy.Tracing.Collector = cfg.Tracing.Collector
	}

	if cfg.TokenManager.JWTSecret != "" {
		cfg.Proxy.TokenManager.JWTSecret = cfg.TokenManager.JWTSecret
	}

	return cfg.Proxy
}

func init() {
	register.AddCommand(ProxyCommand)
}
