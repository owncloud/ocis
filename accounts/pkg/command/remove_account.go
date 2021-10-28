package command

import (
	"fmt"
	"os"

	"github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/flagset"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/urfave/cli/v2"
)

// RemoveAccount command deletes an existing account.
func RemoveAccount(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "Removes an existing account",
		ArgsUsage: "id",
		Aliases:   []string{"rm"},
		Flags:     flagset.RemoveAccountWithConfig(cfg),
		Action: func(c *cli.Context) error {
			accServiceID := cfg.GRPC.Namespace + "." + cfg.Server.Name
			if c.NArg() != 1 {
				fmt.Println("Please provide a user-id")
				os.Exit(1)
			}

			uid := c.Args().First()
			accSvc := accounts.NewAccountsService(accServiceID, grpc.NewClient())
			_, err := accSvc.DeleteAccount(c.Context, &accounts.DeleteAccountRequest{Id: uid})

			if err != nil {
				fmt.Println(fmt.Errorf("could not delete account %w", err))
				return err
			}

			return nil
		}}
}
