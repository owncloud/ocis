package command

import (
	phoenix "github.com/owncloud/reva-phoenix/pkg/command"
	"github.com/spf13/cobra"
)

// Phoenix is the entry point for the phoenix command.
func Phoenix() *cobra.Command {
	cmd := phoenix.Server()
	cmd.Use = "phoenix"
	cmd.Short = "Start phoenix server"

	return cmd
}
