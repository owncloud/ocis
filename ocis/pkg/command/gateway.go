package command

import (
	"github.com/owncloud/ocis/gateway/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GatewayCommand is the entrypoint for the gateway command.
func GatewayCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Gateway.Service.Name,
		Usage:    subcommandDescription(cfg.Gateway.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.Gateway),
	}
}

func init() {
	register.AddCommand(GatewayCommand)
}
