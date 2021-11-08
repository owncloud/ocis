//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/proxy/pkg/command"
	"github.com/urfave/cli/v2"
)

// ProxyCommand is the entry point for the proxy command.
func ProxyCommand(cfg *config.Config) *cli.Command {
	var globalLog shared.Log

	return &cli.Command{
		Name:     "proxy",
		Usage:    "Start proxy server",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Proxy),
		},
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}
			globalLog = cfg.Log
			return nil
		},
		Action: func(c *cli.Context) error {
			// if proxy logging is empty in ocis.yaml
			if (cfg.Proxy.Log == shared.Log{}) && (globalLog != shared.Log{}) {
				// we can safely inherit the global logging values.
				cfg.Proxy.Log = globalLog
			}
			origCmd := command.Server(cfg.Proxy)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(ProxyCommand)
}
