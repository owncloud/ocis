package command

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

// UnifiedRoles bundles available commands for unified roles
func UnifiedRoles(cfg *config.Config) cli.Commands {
	cmds := cli.Commands{
		unifiedRolesStatus(cfg),
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
func unifiedRolesStatus(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list available unified roles",
		Action: func(c *cli.Context) error {
			re := lipgloss.NewRenderer(os.Stdout)
			baseStyle := re.NewStyle().Padding(0, 1)

			var data [][]string

			for _, definition := range unifiedrole.GetBuiltinRoleDefinitionList() {
				data = append(data, []string{"", definition.GetId(), definition.GetDescription()})
			}

			t := table.New().
				Border(lipgloss.NormalBorder()).
				Headers("Enabled", "UID", "Description").
				Rows(data...).
				StyleFunc(func(row, col int) lipgloss.Style {
					if row == 0 {
						return baseStyle.Foreground(lipgloss.Color("252")).Bold(true)
					}

					if row != 0 && col == 0 {
						indicatorStyle := baseStyle.Align(lipgloss.Center)

						// Check if the role is enabled, header takes up the first row
						switch slices.Contains(cfg.UnifiedRoles.AvailableRoles, data[row-1][1]) {
						case true:
							return indicatorStyle.Background(lipgloss.Color("34")) // ANSI green
						default:
							return indicatorStyle.Background(lipgloss.Color("9")) // ANSI red
						}
					}

					return baseStyle
				})

			fmt.Println(t)

			return nil
		},
	}
}
