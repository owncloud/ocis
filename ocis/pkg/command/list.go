package command

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// ListCommand is the entrypoint for the list command.
func ListCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "list",
		Usage:    "Lists running ocis extensions",
		Category: "Runtime",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "hostname",
				Value:       "localhost",
				EnvVars:     []string{"OCIS_RUNTIME_HOST"},
				Destination: &cfg.Runtime.Host,
			},
			&cli.StringFlag{
				Name:        "port",
				Value:       "9250",
				EnvVars:     []string{"OCIS_RUNTIME_PORT"},
				Destination: &cfg.Runtime.Port,
			},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", net.JoinHostPort(cfg.Runtime.Host, cfg.Runtime.Port))
			if err != nil {
				log.Fatal("dialing:", err)
			}

			var arg1 string

			if err := client.Call("Service.List", struct{}{}, &arg1); err != nil {
				log.Fatal(err)
			}

			fmt.Println(arg1)

			return nil
		},
	}
}

func init() {
	register.AddCommand(ListCommand)
}
