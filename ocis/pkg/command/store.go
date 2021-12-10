//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/store/pkg/command"
	"github.com/urfave/cli/v2"
)

// StoreCommand is the entrypoint for the ocs command.
func StoreCommand(cfg *config.Config) *cli.Command {
	var globalLog shared.Log

	return &cli.Command{
		Name:     "store",
		Usage:    "Start a go-micro store",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Store),
		},
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			globalLog = cfg.Log

			return nil
		},
		Action: func(c *cli.Context) error {
			// if accounts logging is empty in ocis.yaml
			if (cfg.Store.Log == shared.Log{}) && (globalLog != shared.Log{}) {
				// we can safely inherit the global logging values.
				cfg.Store.Log = globalLog
			}

			origCmd := command.Server(cfg.Store)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StoreCommand)
}
