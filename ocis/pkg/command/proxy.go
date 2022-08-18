package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/command"
	"github.com/urfave/cli/v2"
)

// ProxyCommand is the entry point for the proxy command.
func ProxyCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Proxy.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Proxy.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.LogError(parser.ParseConfig(cfg, true))
			cfg.Proxy.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Proxy),
	}
}

func init() {
	register.AddCommand(ProxyCommand)
}
