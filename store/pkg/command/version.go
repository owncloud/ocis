package command

import (
	"fmt"
	"os"

	"github.com/owncloud/ocis/ocis-pkg/registry"

	tw "github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/store/pkg/config"
	"github.com/owncloud/ocis/store/pkg/flagset"
	"github.com/urfave/cli/v2"
)

// PrintVersion prints the service versions of all running instances.
func PrintVersion(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print the versions of the running instances",
		Flags: flagset.ListStoreWithConfig(cfg),
		Action: func(c *cli.Context) error {
			reg := registry.GetRegistry()
			services, err := reg.GetService(cfg.Service.Namespace + "." + cfg.Service.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("could not get store services from the registry: %v", err))
				return err
			}

			if len(services) == 0 {
				fmt.Println("No running store service found.")
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
