package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/experimental/pkg/command"
	"github.com/urfave/cli/v2"
)

// ExperimentalCommand is the entrypoint for the experimental command.
func ExperimentalCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Experimental.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Experimental.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Experimental.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Experimental),
	}
}

func init() {
	register.AddCommand(ExperimentalCommand)
}
