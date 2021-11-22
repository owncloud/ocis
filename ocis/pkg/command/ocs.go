//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocs/pkg/command"
	"github.com/urfave/cli/v2"
)

// OCSCommand is the entrypoint for the ocs command.
func OCSCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "ocs",
		Usage:    "Start ocs server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if cfg.Commons != nil {
				cfg.OCS.Commons = cfg.Commons
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(cfg.OCS)
			return handleOriginalAction(c, origCmd)
		},
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.OCS),
		},
	}
}

func init() {
	register.AddCommand(OCSCommand)
}
