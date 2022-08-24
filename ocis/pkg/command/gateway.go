package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/command"
	"github.com/urfave/cli/v2"
)

// GatewayCommand is the entrypoint for the Gateway command.
func GatewayCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Gateway.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Gateway.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.Gateway.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Gateway),
	}
}

func init() {
	register.AddCommand(GatewayCommand)
}
