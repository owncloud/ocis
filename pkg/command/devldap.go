package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-devldap/pkg/command"
	svcconfig "github.com/owncloud/ocis-devldap/pkg/config"
	"github.com/owncloud/ocis-devldap/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// DevLDAPCommand is the entrypoint for the devldap command.
func DevLDAPCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "devldap",
		Usage:    "Start devldap server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.DevLDAP),
		Action: func(c *cli.Context) error {
			scfg := configureDevLDAP(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureDevLDAP(cfg *config.Config) *svcconfig.Config {
	cfg.DevLDAP.Log.Level = cfg.Log.Level
	cfg.DevLDAP.Log.Pretty = cfg.Log.Pretty
	cfg.DevLDAP.Log.Color = cfg.Log.Color
	cfg.DevLDAP.Tracing.Enabled = false
	cfg.DevLDAP.LDAP.Addr = "localhost:9125"
	return cfg.DevLDAP
}

func init() {
	register.AddCommand(DevLDAPCommand)
}
