package command

import (
	"fmt"
	"net/http"

	"github.com/micro/cli"
	"github.com/micro/go-micro/util/log"
	"github.com/owncloud/ocis-graph/pkg/config"
	"github.com/owncloud/ocis-graph/pkg/flagset"
)

// Health is the entrypoint for the health command.
func Health(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "health",
		Usage: "Check health status",
		Flags: flagset.HealthWithConfig(cfg),
		Action: func(c *cli.Context) error {
			resp, err := http.Get(
				fmt.Sprintf(
					"http://%s/healthz",
					cfg.Debug.Addr,
				),
			)

			if err != nil {
				log.Fatalf("Failed to request health check: %w", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Fatalf("Health check responds with [%d]", resp.StatusCode)
			}

			log.Debugf("Health got good state with [%d]", resp.StatusCode)

			return nil
		},
	}
}
