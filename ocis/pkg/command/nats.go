package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/nats/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// NatsCommand is the entrypoint for the Nats command.
func NatsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Nats.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Nats.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Nats.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Nats),
	}
}

func init() {
	register.AddCommand(NatsCommand)
}
