package command

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	ocisinit "github.com/owncloud/ocis/ocis/pkg/init"
	"github.com/owncloud/ocis/ocis/pkg/register"
	cli "github.com/urfave/cli/v2"
)

// InitCommand is the entrypoint for the init command
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
			&cli.BoolFlag{
				Name:    "force-overwrite",
				Aliases: []string{"f"},
				EnvVars: []string{"OCIS_FORCE_CONFIG_OVERWRITE"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:  "config-path",
				Value: defaults.BaseConfigPath(),
				Usage: "config path for the ocis runtime",
			},
		},
		Action: func(c *cli.Context) error {
			insecureFlag := c.String("insecure")
			insecure := false
			if insecureFlag == "ask" {
				answer := strings.ToLower(stringPrompt("Insecure Backends? [Yes|No]"))
				if answer == "yes" || answer == "y" {
					insecure = true
				}
			} else if insecureFlag == "true" {
				insecure = true
			}
			err := ocisinit.CreateConfig(insecure, c.Bool("force-overwrite"), c.String("config-path"))
			if err != nil {
				log.Fatalf("Could not create config: %s", err)
			}
			return nil
		},
	}
}

func init() {
	register.AddCommand(InitCommand)
}

func stringPrompt(label string) string {
	input := ""
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		input, _ = reader.ReadString('\n')
		if input != "" {
			break
		}
	}
	return strings.TrimSpace(input)
}
