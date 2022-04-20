package command

import (
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ocis-config-file",
				Value:       cfg.ConfigFile,
				Usage:       "oCIS config file to be loaded by the extension",
				Destination: &cfg.ConfigFile,
			},
		},
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.Audit),
	}
}

func init() {
	register.AddCommand(AuditCommand)
}
