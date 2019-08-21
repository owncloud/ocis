package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	revaphoenix "github.com/owncloud/reva-phoenix/service"
)

func main() {

	revahyperCommand, allCommandFns := NewRevahHyperCommand()

	basename := filepath.Base(os.Args[0])
	if err := commandFor(basename, revahyperCommand, allCommandFns).Execute(); err != nil {
		os.Exit(1)
	}
}

// NewRevahHyperCommand is the entry point for reva-hyper
func NewRevahHyperCommand() (*cobra.Command, []func() *cobra.Command) {

	apiserver := func() *cobra.Command { return revaphoenix.NewRevaPhoenixCommand("phoenix") }

	commandFns := []func() *cobra.Command{
		apiserver,
	}


	cmd := &cobra.Command{
		Use:   "reva-hyper",
		Short: "Manage oCIS stack",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.Help()
				os.Exit(1)
			}

		},
	}

	for i := range commandFns {
		cmd.AddCommand(commandFns[i]())
	}

	return cmd, commandFns
}

func commandFor(basename string, defaultCommand *cobra.Command, commands []func() *cobra.Command) *cobra.Command {
	for _, commandFn := range commands {
		command := commandFn()
		if command.Name() == basename {
			return command
		}
		for _, alias := range command.Aliases {
			if alias == basename {
				return command
			}
		}
	}

	return defaultCommand
}
