package command

import (
	"fmt"
	"os"

	tw "github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/registry"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
	mreg "go-micro.dev/v4/registry"
)

// VersionCommand is the entrypoint for the version command.
func VersionCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "version",
		Usage:    "print the version of this binary and all running extension instances",
		Category: "info",
		Action: func(c *cli.Context) error {
			fmt.Println("Version: " + version.String)
			fmt.Printf("Compiled: %s\n", version.Compiled())
			fmt.Println("")

			reg := registry.GetRegistry()
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
