package command

import (
	"fmt"

	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/flagset"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
)

// AddAccount command creates a new account
func AddAccount(cfg *config.Config) *cli.Command {
	a := &accounts.Account{
		PasswordProfile: &accounts.PasswordProfile{},
	}
	return &cli.Command{
		Name:    "add",
		Usage:   "Create a new account",
		Aliases: []string{"create", "a"},
		Flags:   flagset.AddAccountWithConfig(cfg, a),
		Before: func(c *cli.Context) error {
			// Write value of username to the flags beneath, as preferred name
			// and on-premises-sam-account-name is probably confusing for users.
			if username := c.String("username"); username != "" {
				if !c.IsSet("on-premises-sam-account-name") {
					if err := c.Set("on-premises-sam-account-name", username); err != nil {
						return err
					}
				}

				if !c.IsSet("preferred-name") {
					if err := c.Set("preferred-name", username); err != nil {
						return err
					}
				}
			}

			return nil

		},
		Action: func(c *cli.Context) error {
			accSvcID := cfg.GRPC.Namespace + "." + cfg.Server.Name
			accSvc := accounts.NewAccountsService(accSvcID, grpc.NewClient())
			_, err := accSvc.CreateAccount(c.Context, &accounts.CreateAccountRequest{
				Account: a,
			})

			if err != nil {
				fmt.Println(fmt.Errorf("could not create account %w", err))
				return err
			}

			return nil
		}}
}
