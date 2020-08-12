// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-proxy/pkg/command"
	svcconfig "github.com/owncloud/ocis-proxy/pkg/config"
	"github.com/owncloud/ocis-proxy/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// ProxyCommand is the entry point for the proxy command.
func ProxyCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "proxy",
		Usage:    "Start proxy server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Proxy),
		Action: func(ctx *cli.Context) error {
			proxyCommand := command.Server(configureProxy(cfg))

			if err := proxyCommand.Before(ctx); err != nil {
				return err
			}

			return cli.HandleAction(proxyCommand.Action, ctx)
		},
	}
}

func configureProxy(cfg *config.Config) *svcconfig.Config {
	cfg.Proxy.Log.Level = cfg.Log.Level
	cfg.Proxy.Log.Pretty = cfg.Log.Pretty
	cfg.Proxy.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Proxy.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Proxy.Tracing.Type = cfg.Tracing.Type
		cfg.Proxy.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Proxy.Tracing.Collector = cfg.Tracing.Collector
		cfg.Proxy.Tracing.Service = cfg.Tracing.Service
	}

	if cfg.Reva.Reva.JWTSecret != "" {
		cfg.Proxy.TokenManager.JWTSecret = cfg.Reva.Reva.JWTSecret
	}

	return cfg.Proxy
}

func init() {
	register.AddCommand(ProxyCommand)
}
