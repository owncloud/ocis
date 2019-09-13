package command

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/oklog/run"
	"github.com/owncloud/ocis-phoenix/pkg/router/debug"
	"github.com/owncloud/ocis-phoenix/pkg/router/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Server is the entrypoint for the server command.
func Server() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start integrated server",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			var gr run.Group

			{
				server := &http.Server{
					Addr: viper.GetString("debug.addr"),
					Handler: debug.Router(
						debug.WithToken(viper.GetString("debug.token")),
						debug.WithPprof(viper.GetBool("debug.pprof")),
					),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", viper.GetString("debug.addr")).
						Msg("Starting debug server")

					if strings.HasPrefix(viper.GetString("debug.addr"), "unix://") {
						socket := strings.TrimPrefix(viper.GetString("debug.addr"), "unix://")

						if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
							log.Error().
								Err(err).
								Str("socket", socket).
								Msg("Failed to remove existing debug socket")

							return err
						}

						listener, err := net.ListenUnix(
							"unix",
							&net.UnixAddr{
								Name: socket,
								Net:  "unix",
							},
						)

						if err != nil {
							log.Error().
								Err(err).
								Msg("Failed to initialize debug unix socket")

							return err
						}

						if err = os.Chmod(socket, os.FileMode(0666)); err != nil {
							log.Error().
								Err(err).
								Msg("Failed to change debug socket permissions")

							return err
						}

						return server.Serve(listener)
					}

					return server.ListenAndServe()
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Error().
							Err(err).
							Msg("Failed to shutdown debug server gracefully")

						return
					}

					if strings.HasPrefix(viper.GetString("debug.addr"), "unix://") {
						socket := strings.TrimPrefix(viper.GetString("debug.addr"), "unix://")

						if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
							log.Error().
								Err(err).
								Str("socket", socket).
								Msg("Failed to remove debug server socket")
						}
					}

					log.Info().
						Err(reason).
						Msg("Shutdown debug server gracefully")
				})
			}

			{
				server := &http.Server{
					Addr: viper.GetString("http.addr"),
					Handler: server.Router(
						server.WithRoot(viper.GetString("http.root")),
						server.WithPath(viper.GetString("asset.path")),
						server.WithConfig(viper.GetString("config.file")),
					),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", viper.GetString("http.addr")).
						Msg("Starting http server")

					if strings.HasPrefix(viper.GetString("http.addr"), "unix://") {
						socket := strings.TrimPrefix(viper.GetString("http.addr"), "unix://")

						if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
							log.Error().
								Err(err).
								Str("socket", socket).
								Msg("Failed to remove existing http socket")

							return err
						}

						listener, err := net.ListenUnix(
							"unix",
							&net.UnixAddr{
								Name: socket,
								Net:  "unix",
							},
						)

						if err != nil {
							log.Error().
								Err(err).
								Msg("Failed to initialize http unix socket")

							return err
						}

						if err = os.Chmod(socket, os.FileMode(0666)); err != nil {
							log.Error().
								Err(err).
								Msg("Failed to change http socket permissions")

							return err
						}

						return server.Serve(listener)
					}

					return server.ListenAndServe()
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Error().
							Err(err).
							Msg("Failed to shutdown http server gracefully")

						return
					}

					if strings.HasPrefix(viper.GetString("http.addr"), "unix://") {
						socket := strings.TrimPrefix(viper.GetString("http.addr"), "unix://")

						if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
							log.Error().
								Err(err).
								Str("socket", socket).
								Msg("Failed to remove http server socket")
						}
					}

					log.Info().
						Err(reason).
						Msg("Shutdown http server gracefully")
				})
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
				})
			}

			return gr.Run()
		},
	}

	cmd.Flags().String("debug-addr", "", "Address to bind debug server")
	viper.BindPFlag("debug.addr", cmd.Flags().Lookup("debug-addr"))
	viper.BindEnv("debug.addr", "PHOENIX_DEBUG_ADDR")

	cmd.Flags().String("debug-token", "", "Token to grant metrics access")
	viper.BindPFlag("debug.token", cmd.Flags().Lookup("debug-token"))
	viper.BindEnv("debug.token", "PHOENIX_DEBUG_TOKEN")

	cmd.Flags().Bool("debug-pprof", false, "Enable pprof debugging")
	viper.BindPFlag("debug.pprof", cmd.Flags().Lookup("debug-pprof"))
	viper.BindEnv("debug.pprof", "PHOENIX_DEBUG_PPROF")

	cmd.Flags().String("http-addr", "", "Address to bind http server")
	viper.BindPFlag("http.addr", cmd.Flags().Lookup("http-addr"))
	viper.BindEnv("http.addr", "PHOENIX_HTTP_ADDR")

	cmd.Flags().String("http-root", "", "Root path for http endpoint")
	viper.BindPFlag("http.root", cmd.Flags().Lookup("http-root"))
	viper.BindEnv("http.root", "PHOENIX_HTTP_ROOT")

	cmd.Flags().String("asset-path", "", "Path to custom assets")
	viper.BindPFlag("asset.path", cmd.Flags().Lookup("asset-path"))
	viper.BindEnv("asset.path", "PHOENIX_ASSET_PATH")

	cmd.Flags().String("config-file", "", "Path to phoenix config")
	viper.BindPFlag("config.file", cmd.Flags().Lookup("config-file"))
	viper.BindEnv("config.file", "PHOENIX_CONFIG_FILE")

	return cmd
}
