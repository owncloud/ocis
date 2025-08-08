package command

import (
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output-format",
				Usage:   "Adjust the basic table output. Available Options: md, colorized",
				Aliases: []string{"o"},
				Value:   "",
			},
		},
		Action: func(c *cli.Context) error {
			headers := []string{"#", "Label", "UID", "Enabled", "Description", "Condition", "Allowed resource actions"}

			opts := []tablewriter.Option{}
			switch c.String("output-format") {
			case "md":
				// Markdown format
				opts = append(opts, tablewriter.WithRenderer(renderer.NewMarkdown()))
			case "colorized":
				opts = append(opts, tablewriter.WithRenderer(renderer.NewColorized()))
			}

			tbl := tablewriter.NewTable(os.Stdout, opts...)
			tbl.Header(headers)

			for i, definition := range unifiedrole.GetRoles(unifiedrole.RoleFilterAll()) {
				const enabled = "enabled"
				const disabled = "disabled"

				rows := [][]string{
					{strconv.Itoa(i + 1), unifiedrole.GetUnifiedRoleLabel(definition.GetId()), definition.GetId(), disabled, definition.GetDescription()},
				}
				if slices.Contains(cfg.UnifiedRoles.AvailableRoles, definition.GetId()) {
					rows[0][3] = enabled
				}

				for i, rolePermission := range definition.GetRolePermissions() {
					actions := strings.Join(rolePermission.GetAllowedResourceActions(), "\n")
					row := []string{rolePermission.GetCondition(), actions}
					switch i {
					case 0:
						rows[0] = append(rows[0], row...)
					default:
						rows[0][5] = rows[0][5] + "\n" + rolePermission.GetCondition()
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
