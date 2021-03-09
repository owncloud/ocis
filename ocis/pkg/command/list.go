package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
)

// ListCommand is the entrypoint for the list command.
func ListCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "list",
		Usage:    "Lists running ocis extensions",
		Category: "Runtime",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "hostname",
				Value:   "localhost",
				EnvVars: []string{"OCIS_RUNTIME_HOSTNAME"},
			},
			&cli.StringFlag{
				Name:    "port",
				Value:   "10666",
				EnvVars: []string{"OCIS_RUNTIME_PORT"},
			},
		},
		Action: func(c *cli.Context) error {
			//client, err := rpc.DialHTTP("tcp", net.JoinHostPort(cfg.Runtime.Hostname, cfg.Runtime.Port))
			//if err != nil {
			//	log.Fatal("dialing:", err)
			//}
			//
			//var arg1 string
			//
			//if err := client.Call("Service.List", struct{}{}, &arg1); err != nil {
			//	log.Fatal(err)
			//}
			//
			//fmt.Println(arg1)

			return nil
		},
	}
}

func init() {
	register.AddCommand(ListCommand)
}
