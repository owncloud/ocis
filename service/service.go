package service

import (
	"os"
	"log"

	"github.com/spf13/cobra"
)

// NewRevaPhoenixCommand is the entry point for reva-phoenix
func NewRevaPhoenixCommand(name string) (*cobra.Command) {
	cmd := &cobra.Command{
		Use:	name,
		Short:	"Start phoenix server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.Help()
				os.Exit(1)
			}else{
				log.Print("Dummy phoenix startup")
			}
		},
	}

	return cmd
}