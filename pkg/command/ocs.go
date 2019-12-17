package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-ocs/pkg/command"
	svcconfig "github.com/owncloud/ocis-ocs/pkg/config"
	"github.com/owncloud/ocis-ocs/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// OCSCommand is the entrypoint for the ocs command.
func OCSCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "ocs",
		Usage:    "Start ocs server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.OCS),
		Action: func(c *cli.Context) error {
			scfg := configureOCS(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureOCS(cfg *config.Config) *svcconfig.Config {
	cfg.OCS.Log.Level = cfg.Log.Level
	cfg.OCS.Log.Pretty = cfg.Log.Pretty
	cfg.OCS.Log.Color = cfg.Log.Color
	cfg.OCS.Tracing.Enabled = false
	cfg.OCS.HTTP.Addr = "localhost:9110"
	cfg.OCS.HTTP.Root = "/"

	return cfg.OCS
}

func init() {
	register.AddCommand(OCSCommand)
}
