package command

import (
	"fmt"
	"os"
	"strconv"

	"github.com/owncloud/ocis/accounts/pkg/flagset"

	"github.com/asim/go-micro/plugins/client/grpc/v4"
	tw "github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/accounts/pkg/config"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/urfave/cli/v2"
)

// ListAccounts command lists all accounts
func ListAccounts(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List existing accounts",
		Aliases: []string{"ls"},
		Flags:   flagset.Root(cfg),
		Action: func(c *cli.Context) error {
			accSvcID := cfg.GRPC.Namespace + "." + cfg.Server.Name
			accSvc := accounts.NewAccountsService(accSvcID, grpc.NewClient())
			resp, err := accSvc.ListAccounts(c.Context, &accounts.ListAccountsRequest{})

			if err != nil {
				fmt.Println(fmt.Errorf("could not list accounts %w", err))
				return err
			}

			buildAccountsListTable(resp.Accounts).Render()
			return nil
		}}
}

// buildAccountsListTable creates an ascii table for printing on the cli
func buildAccountsListTable(accs []*accounts.Account) *tw.Table {
	table := tw.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "DisplayName", "Mail", "AccountEnabled"})
	table.SetAutoFormatHeaders(false)
	for _, acc := range accs {
		table.Append([]string{
			acc.Id,
			acc.DisplayName,
			acc.Mail,
			strconv.FormatBool(acc.AccountEnabled)})
	}
	return table
}
