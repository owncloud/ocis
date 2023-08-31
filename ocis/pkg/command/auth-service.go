package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/auth-service/pkg/command"
	"github.com/urfave/cli/v2"
)

// AuthServiceCommand is the entrypoint for the AuthService command.
func AuthServiceCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AuthService.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.AuthService.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.AuthService.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.AuthService),
	}
}

func init() {
	register.AddCommand(AuthServiceCommand)
}
