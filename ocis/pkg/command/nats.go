package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/nats/pkg/command"
	"github.com/urfave/cli/v2"
)

// NatsCommand is the entrypoint for the Nats command.
func NatsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Nats.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Nats.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.Nats.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Nats),
	}
}

func init() {
	register.AddCommand(NatsCommand)
}
