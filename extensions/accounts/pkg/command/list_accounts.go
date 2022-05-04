package command

import (
	"fmt"
	"os"
	"strconv"

	accountsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/accounts/v0"
	accountssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/accounts/v0"

	"github.com/owncloud/ocis/v2/extensions/accounts/pkg/flagset"

	"github.com/go-micro/plugins/v4/client/grpc"
	tw "github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/v2/extensions/accounts/pkg/config"
	"github.com/urfave/cli/v2"
)

// ListAccounts command lists all accounts
func ListAccounts(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "list",
		Usage:    "list existing accounts",
		Category: "account management",
		Aliases:  []string{"ls"},
		Flags:    flagset.ListAccountsWithConfig(cfg),
		Action: func(c *cli.Context) error {
			accSvcID := cfg.GRPC.Namespace + "." + cfg.Service.Name
			accSvc := accountssvc.NewAccountsService(accSvcID, grpc.NewClient())
			resp, err := accSvc.ListAccounts(c.Context, &accountssvc.ListAccountsRequest{})

			if err != nil {
				fmt.Println(fmt.Errorf("could not list accounts %w", err))
				return err
			}

			buildAccountsListTable(resp.Accounts).Render()
			return nil
		}}
}

// buildAccountsListTable creates an ascii table for printing on the cli
func buildAccountsListTable(accs []*accountsmsg.Account) *tw.Table {
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
