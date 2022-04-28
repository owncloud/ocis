package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/audit/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AuditCommand is the entrypoint for the audit command.
func AuditCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "audit",
		Usage:    "start audit service",
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
			}
			return err
		},
		Subcommands: command.GetCommands(cfg.Audit),
	}
}

func init() {
	register.AddCommand(AuditCommand)
}
