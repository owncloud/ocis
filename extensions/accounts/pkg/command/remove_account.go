package command

import (
	"fmt"
	"os"

	accountssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/accounts/v0"

	"github.com/owncloud/ocis/v2/extensions/accounts/pkg/flagset"

	"github.com/go-micro/plugins/v4/client/grpc"
	"github.com/owncloud/ocis/v2/extensions/accounts/pkg/config"
	"github.com/urfave/cli/v2"
)

// RemoveAccount command deletes an existing account.
func RemoveAccount(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "removes an existing account",
		Category:  "account management",
		ArgsUsage: "id",
		Aliases:   []string{"rm"},
		Flags:     flagset.RemoveAccountWithConfig(cfg),
		Action: func(c *cli.Context) error {
			accServiceID := cfg.GRPC.Namespace + "." + cfg.Service.Name
			if c.NArg() != 1 {
				fmt.Println("Please provide a user-id")
				os.Exit(1)
			}

			uid := c.Args().First()
			accSvc := accountssvc.NewAccountsService(accServiceID, grpc.NewClient())
			_, err := accSvc.DeleteAccount(c.Context, &accountssvc.DeleteAccountRequest{Id: uid})

			if err != nil {
				fmt.Println(fmt.Errorf("could not delete account %w", err))
				return err
			}

			return nil
		}}
}
