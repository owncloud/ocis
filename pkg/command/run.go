package command

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
	"github.com/refs/pman/pkg/process"
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
				Value:       "10666",
				EnvVars:     []string{"OCIS_RUNTIME_PORT"},
				Destination: &cfg.Runtime.Port,
			},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", net.JoinHostPort(cfg.Runtime.Hostname, cfg.Runtime.Port))
			if err != nil {
				log.Fatal("dialing:", err)
			}

			proc := process.NewProcEntry(os.Args[2], []string{os.Args[2]}...)
			var res int

			if err := client.Call("Service.Start", proc, &res); err != nil {
				log.Fatal(err)
			}

			fmt.Println(res)
			return nil
		},
	}
}

func init() {
	register.AddCommand(RunCommand)
}
