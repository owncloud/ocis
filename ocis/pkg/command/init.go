package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	cli "github.com/urfave/cli/v2"
)

func InitCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialise an ocis config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "insecure",
				EnvVars: []string{"OCIS_INSECURE"},
				Value:   "ask",
			},
		},
		Action: func(c *cli.Context) error {
			// TODO: discuss if we want overwrite protection for existing configs
			insecureFlag := c.String("insecure")
			if insecureFlag == "ask" {
				answer := strings.ToLower(StringPrompt("Insecure Backends? [Yes|No]"))
				if answer == "yes" || answer == "y" {
					cfg.Proxy.InsecureBackends = true
				} else {
					cfg.Proxy.InsecureBackends = false
				}
			} else {
				if insecureFlag == "true" {
					cfg.Proxy.InsecureBackends = true
				} else {
					cfg.Proxy.InsecureBackends = false
				}
			}
			fmt.Println(cfg.Proxy.InsecureBackends)
			return nil
		},
	}
}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func init() {
	register.AddCommand(InitCommand)
}
