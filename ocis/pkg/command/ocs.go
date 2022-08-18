package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/command"
	"github.com/urfave/cli/v2"
)

// OCSCommand is the entrypoint for the ocs command.
func OCSCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.OCS.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.OCS.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.LogError(parser.ParseConfig(cfg, true))
			cfg.OCS.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.OCS),
	}
}

func init() {
	register.AddCommand(OCSCommand)
}
