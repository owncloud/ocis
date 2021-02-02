package cmd

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/spf13/cobra"
)

// Kill an extension.
func Kill(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "kill",
		Aliases: []string{"k"},
		Short:   "Kill a running extensions.",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := rpc.DialHTTP("tcp", net.JoinHostPort(cfg.Hostname, cfg.Port))
			if err != nil {
				log.Fatal("dialing:", err)
			}

			var arg1 int

			if err := client.Call("Service.Kill", &args[0], &arg1); err != nil {
				log.Fatal(err)
			}

			fmt.Println(arg1)
		},
	}
}
