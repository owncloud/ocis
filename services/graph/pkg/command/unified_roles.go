package command

import (
	"os"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

// UnifiedRoles bundles available commands for unified roles
func UnifiedRoles(cfg *config.Config) cli.Commands {
	cmds := cli.Commands{
		listUnifiedRoles(cfg),
	}

	for _, cmd := range cmds {
		cmd.Category = "unified-roles"
		cmd.Name = strings.Join([]string{cmd.Name, "unified-roles"}, "-")
		cmd.Before = func(c *cli.Context) error {
			return configlog.ReturnError(parser.ParseConfig(cfg))
		}
	}

	return cmds
}

// unifiedRolesStatus lists available unified roles, it contains an indicator to show if the role is enabled or not
func listUnifiedRoles(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list available unified roles",
		Action: func(c *cli.Context) error {
			tbl := tablewriter.NewWriter(os.Stdout)
			tbl.SetRowLine(true)
			tbl.SetAutoMergeCellsByColumnIndex([]int{0}) // rowspan should only affect the first column

			headers := []string{"UID", "Enabled", "Description", "Condition", "Allowed resource actions"}
			tbl.SetHeader(headers)

			for _, definition := range unifiedrole.GetRoles(unifiedrole.RoleFilterAll()) {
				const enabled = "enabled"
				const disabled = "disabled"

				rows := [][]string{
					{definition.GetId(), disabled, definition.GetDescription()},
				}
				if slices.Contains(cfg.UnifiedRoles.AvailableRoles, definition.GetId()) {
					rows[0][1] = enabled
				}

				for i, rolePermission := range definition.GetRolePermissions() {
					actions := strings.Join(rolePermission.GetAllowedResourceActions(), "\n")
					row := []string{rolePermission.GetCondition(), actions}
					switch i {
					case 0:
						rows[0] = append(rows[0], row...)
					default:
						rows = append(rows, append(slices.Clone(rows[0][:len(rows[0])-len(row)]), row...))
					}
				}

				for _, row := range rows {
					// balance the row before adding it to the table,
					// this prevents the row from having empty columns.
					tbl.Append(append(row, make([]string, len(headers)-len(row))...))
				}
			}

			tbl.Render()
			return nil
		},
	}
}
