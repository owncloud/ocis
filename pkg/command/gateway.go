package command

import (
	"context"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis-reva/pkg/server/debug"
	"github.com/owncloud/ocis-reva/pkg/service/external"
)

// Gateway is the entrypoint for the gateway command.
func Gateway(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "gateway",
		Usage: "Start reva gateway",
		Flags: flagset.GatewayWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Gateway.Services = c.StringSlice("service")

			return nil
		},
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

			{

				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus": cfg.Reva.Gateway.MaxCPUs,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
						"gatewaysvc": cfg.Reva.Gateway.URL, // Todo or address?
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.Gateway.Network,
						"address": cfg.Reva.Gateway.Addr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"gateway": map[string]interface{}{
								// registries is located on the gateway
								"authregistrysvc":    cfg.Reva.Gateway.URL,
								"storageregistrysvc": cfg.Reva.Gateway.URL,
								"appregistrysvc":     cfg.Reva.Gateway.URL,
								// user metadata is located on the users services
								"preferencessvc":  cfg.Reva.Users.URL,
								"userprovidersvc": cfg.Reva.Users.URL,
								// sharing is located on the sharing service
								"usershareprovidersvc":          cfg.Reva.Sharing.URL,
								"publicshareprovidersvc":        cfg.Reva.Sharing.URL,
								"ocmshareprovidersvc":           cfg.Reva.Sharing.URL,
								"commit_share_to_storage_grant": cfg.Reva.Gateway.CommitShareToStorageGrant,
								"commit_share_to_storage_ref":   cfg.Reva.Gateway.CommitShareToStorageRef,
								"share_folder":                  cfg.Reva.Gateway.ShareFolder, // ShareFolder is the location where to create shares in the recipient's storage provider.
								// other
								"disable_home_creation_on_login": cfg.Reva.Gateway.DisableHomeCreationOnLogin,
								"datagateway":                    cfg.Reva.Frontend.URL,
								"transfer_shared_secret":         cfg.Reva.TransferSecret,
								"transfer_expires":               cfg.Reva.TransferExpires,
							},
							"authregistry": map[string]interface{}{
								"driver": "static",
								"drivers": map[string]interface{}{
									"static": map[string]interface{}{
										"rules": map[string]interface{}{
											"basic":  cfg.Reva.AuthBasic.URL,
											"bearer": cfg.Reva.AuthBearer.URL,
										},
									},
								},
							},
							"storageregistry": map[string]interface{}{
								"driver": "static",
								"drivers": map[string]interface{}{
									"static": map[string]interface{}{
										"home_provider": cfg.Reva.Gateway.HomeProvider,
										"rules": map[string]interface{}{
											cfg.Reva.StorageRoot.MountPath: cfg.Reva.StorageRoot.URL,
											cfg.Reva.StorageRoot.MountID:   cfg.Reva.StorageRoot.URL,
											cfg.Reva.StorageHome.MountPath: cfg.Reva.StorageHome.URL,
											// the home storage has no mount id. In responses it returns the mount id of the actual storage
											cfg.Reva.StorageEOS.MountPath:    cfg.Reva.StorageEOS.URL,
											cfg.Reva.StorageEOS.MountID:      cfg.Reva.StorageEOS.URL,
											cfg.Reva.StorageOC.MountPath:     cfg.Reva.StorageOC.URL,
											cfg.Reva.StorageOC.MountID:       cfg.Reva.StorageOC.URL,
											cfg.Reva.StorageS3.MountPath:     cfg.Reva.StorageS3.URL,
											cfg.Reva.StorageS3.MountID:       cfg.Reva.StorageS3.URL,
											cfg.Reva.StorageWND.MountPath:    cfg.Reva.StorageWND.URL,
											cfg.Reva.StorageWND.MountID:      cfg.Reva.StorageWND.URL,
											cfg.Reva.StorageCustom.MountPath: cfg.Reva.StorageCustom.URL,
											cfg.Reva.StorageCustom.MountID:   cfg.Reva.StorageCustom.URL,
										},
									},
								},
							},
						},
					},
				}

				gr.Add(func() error {
					err := external.RegisterGRPCEndpoint(
						ctx,
						"com.owncloud.reva",
						uuid.String(),
						cfg.Reva.Gateway.Addr,
						logger,
					)

					if err != nil {
						return err
					}

					runtime.Run(rcfg, pidFile)
					return nil
				}, func(_ error) {
					logger.Info().
						Str("server", c.Command.Name).
						Msg("Shutting down server")

					cancel()
				})

			}

			{
				server, err := debug.Server(
					debug.Name(c.Command.Name+"-debug"),
					debug.Addr(cfg.Reva.Gateway.DebugAddr),
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
