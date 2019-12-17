package runtime

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro/config/cmd"
	"github.com/owncloud/ocis-pkg/log"
)

// Command adds micro runtime commands to the cli app
func Command(app *cli.App) cli.Command {
	command := cli.Command{
		Name:        "micro-runtime",
		Description: "starts the go-micro runtime and its services",
		Category:    "Base",
		Action: func(c *cli.Context) error {
			runtime := Runtime{
				Services: RuntimeServices,
				R:        cmd.DefaultCmd.Options().Runtime,
				Logger:   log.NewLogger(),
			}

			{
				runtime.Start()
				runtime.Trap()
			}

			return nil
		},
	}
	return command
}
