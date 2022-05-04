package command

import (
	"fmt"

	accountsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/accounts/v0"
	accountssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/accounts/v0"

	"github.com/go-micro/plugins/v4/client/grpc"
	"github.com/owncloud/ocis/v2/extensions/accounts/pkg/config"
	"github.com/owncloud/ocis/v2/extensions/accounts/pkg/flagset"
	"github.com/urfave/cli/v2"
)

// AddAccount command creates a new account
func AddAccount(cfg *config.Config) *cli.Command {
	a := &accountsmsg.Account{
		PasswordProfile: &accountsmsg.PasswordProfile{},
	}
	return &cli.Command{
		Name:     "add",
		Usage:    "create a new account",
		Category: "account management",
		Aliases:  []string{"create", "a"},
		Flags:    flagset.AddAccountWithConfig(cfg, a),
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
			accSvcID := cfg.GRPC.Namespace + "." + cfg.Service.Name
			accSvc := accountssvc.NewAccountsService(accSvcID, grpc.NewClient())
			_, err := accSvc.CreateAccount(c.Context, &accountssvc.CreateAccountRequest{
				Account: a,
			})

			if err != nil {
				fmt.Println(fmt.Errorf("could not create account %w", err))
				return err
			}

			return nil
		}}
}
