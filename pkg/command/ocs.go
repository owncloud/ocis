package command

import (
	ocs "github.com/owncloud/ocis-ocs/pkg/command"
	"github.com/spf13/cobra"
)

// Ocs is the entry point for the ocs command.
func Ocs() *cobra.Command {
	cmd := ocs.Server()
	cmd.Use = "ocs"
	cmd.Short = "Start ocs server"

	return cmd
}
