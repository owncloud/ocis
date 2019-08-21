package service

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// func main() {

// 	revaphoenixCommand, allCommandFns := NewRevaPhoenixCommand()

// 	basename := filepath.Base(os.Args[0])
// 	if err := commandFor(basename, revaphoenixCommand, allCommandFns).Execute(); err != nil {
// 		os.Exit(1)
// 	}
// }

// NewRevaPhoenixCommand is the entry point for reva-phoenix
func NewRevaPhoenixCommand() (*cobra.Command, []func() *cobra.Command) {

	cmd := &cobra.Command{
		Use:   "reva-phoenix",
		Short: "Request a new project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.Help()
				os.Exit(1)
			}

		},
	}

	return cmd, nil
}

// func commandFor(basename string, defaultCommand *cobra.Command, commands []func() *cobra.Command) *cobra.Command {
// 	for _, commandFn := range commands {
// 		command := commandFn()
// 		if command.Name() == basename {
// 			return command
// 		}
// 		for _, alias := range command.Aliases {
// 			if alias == basename {
// 				return command
// 			}
// 		}
// 	}

// 	return defaultCommand
// }
