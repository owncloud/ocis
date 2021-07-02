package cmd

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/spf13/cobra"
)

// Run an extension.
func Run(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run an extension.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := rpc.DialHTTP("tcp", net.JoinHostPort(cfg.Hostname, cfg.Port))
			if err != nil {
				log.Fatal("dialing:", err)
			}
			var res int
			if err := client.Call("Service.Start", &args[0], &res); err != nil {
				log.Fatal(err)
			}

			fmt.Println(res)
		},
	}
}
