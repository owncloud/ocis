package command

import (
	"fmt"
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// Health is the entrypoint for the health command.
func Health(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "health",
		Usage: "Check health status",
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg)
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			resp, err := http.Get(
				fmt.Sprintf(
					"http://%s/healthz",
					cfg.Debug.Addr,
				),
			)

			if err != nil {
				logger.Fatal().
					Err(err).
					Msg("Failed to request health check")
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Fatal().
					Int("code", resp.StatusCode).
					Msg("Health seems to be in bad state")
			}

			logger.Debug().
				Int("code", resp.StatusCode).
				Msg("Health got a good state")

			return nil
		},
	}
}

func init() {
	register.AddCommand(Health)
}
