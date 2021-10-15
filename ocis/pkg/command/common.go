package command

import "github.com/urfave/cli/v2"

func handleOriginalAction(c *cli.Context, cmd *cli.Command) error {

	if cmd.Before != nil {
		if err := cmd.Before(c); err != nil {
			return err
		}
	}

	return cli.HandleAction(cmd.Action, c)
}
