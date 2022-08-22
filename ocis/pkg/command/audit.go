package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/audit/pkg/command"
	"github.com/urfave/cli/v2"
)

// AuditCommand is the entrypoint for the Audit command.
func AuditCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Audit.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Audit.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.Audit.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Audit),
	}
}

func init() {
	register.AddCommand(AuditCommand)
}
