package runtime

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/owncloud/ocis-pkg/v2/log"
)

// Command adds micro runtime commands to the cli app
func Command(app *cli.App) *cli.Command {
	command := cli.Command{
		Name:        "micro",
		Description: "starts the go-micro runtime services",
		Category:    "Micro",
		Action: func(c *cli.Context) error {
			runtime := New(
				Services(MicroServices),
				Logger(log.NewLogger()),
				MicroRuntime(cmd.DefaultCmd.Options().Runtime),
			)

			{
				runtime.Start()
				runtime.Trap()
			}

			return nil
		},
	}
	return &command
}
