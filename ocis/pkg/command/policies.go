package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/policies/pkg/command"
	"github.com/urfave/cli/v2"
)

// PoliciesCommand is the entrypoint for the policies service.
func PoliciesCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Policies.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Policies.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.Policies.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Policies),
	}
}

func init() {
	register.AddCommand(PoliciesCommand)
}
