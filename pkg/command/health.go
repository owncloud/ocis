package command

import (
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
			log.Info().
				Str("addr", viper.GetString("debug.addr")).
				Msg("Executed health command")
		},
	}

	cmd.Flags().String("debug-addr", "", "Address to debug endpoint")
	viper.BindPFlag("debug.addr", cmd.Flags().Lookup("debug-addr"))
	viper.BindEnv("debug.addr", "PHOENIX_DEBUG_ADDR")

	return cmd
}
