package command

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	cli "github.com/micro/cli/v2"

	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
)

// RunCommand is the entrypoint for the run command.
func RunCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "run",
		Usage:    "Runs an extension",
		Category: "Runtime",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "hostname",
				Value:       "localhost",
				EnvVars:     []string{"OCIS_RUNTIME_HOSTNAME"},
				Destination: &cfg.Runtime.Hostname,
			},
			&cli.StringFlag{
				Name:        "port",
				Value:       "6060",
				EnvVars:     []string{"OCIS_RUNTIME_PORT"},
				Destination: &cfg.Runtime.Port,
			},
		},
		Action: func(c *cli.Context) error {
			// TODO(refs) this implementation changes as we don't depend on os threads anymore.
			client, err := rpc.DialHTTP("tcp", net.JoinHostPort(cfg.Runtime.Hostname, cfg.Runtime.Port))
			if err != nil {
				log.Fatal("dialing:", err)
			}

			var reply int

			if err := client.Call("Service.Start", "settings", &reply); err != nil {
				log.Fatal(err)
			}
			fmt.Println(reply)

			return nil
		},
	}
}

func init() {
	register.AddCommand(RunCommand)
}
