package cmd

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/owncloud/ocis/ocis/pkg/runtime/process"
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

			proc := process.NewProcEntry(args[0], os.Environ(), []string{args[0]}...)
			var res int

			if err := client.Call("Service.Start", proc, &res); err != nil {
				log.Fatal(err)
			}

			fmt.Println(res)
		},
	}
}
