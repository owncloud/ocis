package command

import (
	"context"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis-reva/pkg/server/debug"
)

// AuthProvider is the entrypoint for the authprovider command.
func AuthProvider(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "authprovider",
		Usage: "Start authprovider server",
		Flags: flagset.ServerWithConfig(cfg),
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					logger.Error().
						Str("type", t).
						Msg("Reva only supports the jaeger tracing backend")

				case "jaeger":
					logger.Info().
						Str("type", t).
						Msg("configuring reva to use the jaeger tracing backend")

				case "zipkin":
					logger.Error().
						Str("type", t).
						Msg("Reva only supports the jaeger tracing backend")

				default:
					logger.Warn().
						Str("type", t).
						Msg("Unknown tracing backend")
				}

			} else {
				logger.Debug().
					Msg("Tracing is not enabled")
			}

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
				//metrics     = metrics.New()
			)

			defer cancel()

			// TODO Flags have to be injected all the way down to the go-micro service
			{

				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus": cfg.Reva.MaxCPUs,
					},
					"grpc": map[string]interface{}{
						"network":          cfg.Reva.GRPC.Network,
						"address":          cfg.Reva.GRPC.Addr,
						"enabled_services": []string{"authprovider"},
						"interceptors": map[string]interface{}{
							"auth": map[string]interface{}{
								"token_manager": "jwt",
								"token_managers": map[string]interface{}{
									"jwt": map[string]interface{}{
										"secret": cfg.Reva.JWTSecret,
									},
								},
								"skip_methods": []string{
									// we need to allow calls that happen during authentication
									"/cs3.authproviderv0alpha.AuthProviderService/Authenticate",
									"/cs3.userproviderv0alpha.UserProviderService/GetUser",
								},
							},
						},
						"services": map[string]interface{}{
							"authprovider": map[string]interface{}{
								"auth_manager": "oidc",
								"auth_managers": map[string]interface{}{
									"oidc": map[string]interface{}{
										"provider": cfg.AuthProvider.Provider,
										"insecure": cfg.AuthProvider.Insecure,
									},
								},
							},
						},
					},
				}
				// TODO merge configs for the same address

				gr.Add(func() error {
					runtime.Run(rcfg, pidFile)
					return nil
				}, func(_ error) {
					logger.Info().
						Str("server", "authprovider").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("server", "debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)

					defer timeout()
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						logger.Info().
							Err(err).
							Str("server", "debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("server", "debug").
							Msg("Shutting down server")
					}
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
					cancel()
				})
			}

			return gr.Run()
		},
	}
}
