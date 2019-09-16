package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Health is the entrypoint for the health command.
func Health() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Check health status",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := http.Get(
				fmt.Sprintf(
					"http://%s/healthz",
					viper.GetString("debug.addr"),
				),
			)

			if err != nil {
				log.Error().
					Err(err).
					Msg("Failed to request health check")

				os.Exit(1)
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Error().
					Int("code", resp.StatusCode).
					Msg("Health seems to be in bad state")

				os.Exit(1)
			}

			os.Exit(0)
		},
	}

	cmd.Flags().String("debug-addr", "", "Address to debug endpoint")
	viper.BindPFlag("debug.addr", cmd.Flags().Lookup("debug-addr"))
	viper.BindEnv("debug.addr", "WEBDAV_DEBUG_ADDR")

	return cmd
}
