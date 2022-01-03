package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/store/pkg/command"
	"github.com/urfave/cli/v2"
)

// StoreCommand is the entrypoint for the ocs command.
func StoreCommand(cfg *config.Config) *cli.Command {

	return &cli.Command{
		Name:     "store",
		Usage:    "Start a go-micro store",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			if cfg.Commons != nil {
				cfg.Store.Commons = cfg.Commons
			}

			return nil
		},
		Subcommands: command.GetCommands(cfg.Store),
	}
}

func init() {
	register.AddCommand(StoreCommand)
}
