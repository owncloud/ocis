package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis/pkg/client"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/flagset"
	"github.com/owncloud/ocis/pkg/register"
)

// Login is the entrypoint for the login command.
func Login(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "Login and print the oauth2 token",
		Flags: flagset.LoginWithConfig(cfg),
		Action: func(c *cli.Context) error {
			client.HandleOpenIDFlow(cfg.OIDC)
			return nil
		},
	}
}

func init() {
	register.AddCommand(Login)
}
