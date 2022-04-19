package command

import (
	"github.com/owncloud/ocis/extensions/audit/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AuditCommand is the entrypoint for the audit command.
func AuditCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        cfg.Audit.Service.Name,
		Usage:       subcommandDescription(cfg.Audit.Service.Name),
		Category:    "extensions",
		Subcommands: command.GetCommands(cfg.Audit),
	}
}

func init() {
	register.AddCommand(AuditCommand)
}
