package command

import (
	"github.com/owncloud/ocis/extensions/web/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// WebCommand is the entrypoint for the web command.
func WebCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Web.Service.Name,
		Usage:    subcommandDescription(cfg.Web.Service.Name),
		Category: "extensions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ocis-config-file",
				Value:       cfg.ConfigFile,
				Usage:       "oCIS config file to be loaded by the runtime and extensions",
				Destination: &cfg.ConfigFile,
			},
		},
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.Web),
	}
}

func init() {
	register.AddCommand(WebCommand)
}
