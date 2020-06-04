// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaFrontendCommand is the entrypoint for the reva-frontend command.
func RevaFrontendCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-frontend",
		Usage:    "Start reva frontend",
		Category: "Extensions",
		Flags:    flagset.FrontendWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaFrontend(cfg)

			return cli.HandleAction(
				command.Frontend(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaFrontend(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaFrontendCommand)
}
