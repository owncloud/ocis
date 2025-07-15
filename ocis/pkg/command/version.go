package command

import (
	"fmt"
	"os"

	tw "github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	mreg "go-micro.dev/v4/registry"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
)

const (
	_skipServiceListingFlagName = "skip-services"
)

// VersionCommand is the entrypoint for the version command.
func VersionCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "print the version of this binary and all running service instances",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  _skipServiceListingFlagName,
				Usage: "skip service listing",
			},
		},
		Category: "info",
		Action: func(c *cli.Context) error {
			fmt.Println("Version: " + version.GetString())
			fmt.Printf("Compiled: %s\n", version.Compiled())

			if c.Bool(_skipServiceListingFlagName) {
				return nil
			}

			fmt.Print("\n")

			reg := registry.GetRegistry()
			serviceList, err := reg.ListServices()
			if err != nil {
				fmt.Printf("could not list services: %v\n", err)
				return err
			}

			var services []*mreg.Service
			for _, s := range serviceList {
				s, err := reg.GetService(s.Name)
				if err != nil {
					fmt.Printf("could not get service: %v\n", err)
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
