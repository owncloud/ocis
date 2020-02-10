package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaGatewayCommand is the entrypoint for the reva-gateway command.
func RevaGatewayCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-gateway",
		Usage:    "Start reva gateway",
		Category: "Extensions",
		Flags:    flagset.GatewayWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaGateway(cfg)

			return cli.HandleAction(
				command.Gateway(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaGateway(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaGatewayCommand)
}
