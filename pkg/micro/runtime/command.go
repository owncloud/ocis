package runtime

import (
	"github.com/micro/cli/v2"
)

// Command adds micro runtime commands to the cli app
func Command(app *cli.App) *cli.Command {
	command := cli.Command{
		Name:        "micro",
		Description: "starts the go-micro runtime services",
		Category:    "Micro",
		Action: func(c *cli.Context) error {
			runtime := New()
			runtime.Start()

			return nil
		},
	}
	return &command
}
