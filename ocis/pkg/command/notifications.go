package command

import (
	"github.com/owncloud/ocis/extensions/notifications/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// NatsServerCommand is the entrypoint for the nats server command.
func NotificationsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "notifications",
		Usage:    "start notifications service",
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.Notifications),
	}
}

func init() {
	register.AddCommand(NotificationsCommand)
}
