//go:build !simple
// +build !simple

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
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Store),
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(cfg.Store)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StoreCommand)
}
