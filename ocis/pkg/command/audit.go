package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/audit/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AuditCommand is the entrypoint for the Audit command.
func AuditCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Audit.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Audit.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Audit.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Audit),
	}
}

func init() {
	register.AddCommand(AuditCommand)
}
