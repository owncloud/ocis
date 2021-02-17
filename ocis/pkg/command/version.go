package command

import (
	"fmt"
	"os"

	mreg "github.com/asim/go-micro/v3/registry"
	"github.com/micro/cli/v2"
	tw "github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/ocis-pkg/registry"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
)

// VersionCommand is the entrypoint for the version command.
func VersionCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "version",
		Usage:    "Lists running services with version",
		Category: "Runtime",
		Action: func(c *cli.Context) error {
			reg := *registry.GetRegistry()
			serviceList, err := reg.ListServices()
			if err != nil {
				fmt.Println(fmt.Errorf("could not list services: %v", err))
				return err
			}

			var services []*mreg.Service
			for _, s := range serviceList {
				s, err := reg.GetService(s.Name)
				if err != nil {
					fmt.Println(fmt.Errorf("could not get service: %v", err))
					return err
				}
				services = append(services, s...)
			}

			if len(services) == 0 {
				fmt.Println("No running services found.")
				return nil
			}

			table := tw.NewWriter(os.Stdout)
			table.SetHeader([]string{"Version", "Address", "Id"})
			table.SetAutoFormatHeaders(false)
			for _, s := range services {
				for _, n := range s.Nodes {
					table.Append([]string{s.Version, n.Address, n.Id})
				}
			}
			table.Render()
			return nil
		},
	}
}

func init() {
	register.AddCommand(VersionCommand)
}
