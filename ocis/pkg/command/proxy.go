package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/proxy/pkg/command"
	"github.com/urfave/cli/v2"
)

// ProxyCommand is the entry point for the proxy command.
func ProxyCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Proxy.Service.Name,
		Usage:    subcommandDescription(cfg.Proxy.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.Proxy),
	}
}

func init() {
	register.AddCommand(ProxyCommand)
}
