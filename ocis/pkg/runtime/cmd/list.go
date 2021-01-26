package cmd

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/spf13/cobra"
)

// List running extensions.
func List(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"r"},
		Short:   "List running extensions",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := rpc.DialHTTP("tcp", net.JoinHostPort(cfg.Hostname, cfg.Port))
			if err != nil {
				log.Fatal("dialing:", err)
			}

			var arg1 string

			if err := client.Call("Service.List", struct{}{}, &arg1); err != nil {
				log.Fatal(err)
			}

			fmt.Println(arg1)
		},
	}
}
