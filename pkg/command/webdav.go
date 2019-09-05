package command

import (
	webdav "github.com/owncloud/ocis-webdav/pkg/command"
	"github.com/spf13/cobra"
)

// Webdav is the entry point for the webdav command.
func Webdav() *cobra.Command {
	cmd := webdav.Server()
	cmd.Use = "webdav"
	cmd.Short = "Start webdav server"

	return cmd
}
