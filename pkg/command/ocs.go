// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-ocs/pkg/command"
	svcconfig "github.com/owncloud/ocis-ocs/pkg/config"
	"github.com/owncloud/ocis-ocs/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// OCSCommand is the entrypoint for the ocs command.
func OCSCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "ocs",
		Usage:    "Start ocs server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.OCS),
		Action: func(ctx *cli.Context) error {
			ocsCommand := command.Server(configureOCS(cfg))

			if err := ocsCommand.Before(ctx); err != nil {
				return err
			}

			return cli.HandleAction(ocsCommand.Action, ctx)
		},
	}
}

func configureOCS(cfg *config.Config) *svcconfig.Config {
	cfg.OCS.Log.Level = cfg.Log.Level
	cfg.OCS.Log.Pretty = cfg.Log.Pretty
	cfg.OCS.Log.Color = cfg.Log.Color

	return cfg.OCS
}

func init() {
	register.AddCommand(OCSCommand)
}
