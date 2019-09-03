package service

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// HealthCommand is the entrypoint for the health command.
func HealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Check health status",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().
				Str("addr", viper.GetString("metrics.addr")).
				Msg("Executed health command")
		},
	}

	cmd.Flags().String("metrics-addr", "", "Address to metrics endpoint")
	viper.BindPFlag("metrics.addr", cmd.Flags().Lookup("metrics-addr"))
	viper.BindEnv("metrics.addr", "PHOENIX_METRICS_ADDR")

	return cmd
}
