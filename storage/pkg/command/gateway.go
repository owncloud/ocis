package command

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/service/external"
)

// Gateway is the entrypoint for the gateway command.
func Gateway(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "gateway",
		Usage: "Start gateway",
		Flags: flagset.GatewayWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Gateway.Services = c.StringSlice("service")
			cfg.Reva.StorageRegistry.Rules = c.StringSlice("storage-registry-rule")

			if cfg.Reva.DataGateway.PublicURL == "" {
				cfg.Reva.DataGateway.PublicURL = strings.TrimRight(cfg.Reva.Frontend.PublicURL, "/") + "/data"
			}

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
						Msg("configuring storage to use the jaeger tracing backend")

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
						"max_cpus":             cfg.Reva.Users.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": c.Command.Name,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
						"gatewaysvc": cfg.Reva.Gateway.Endpoint,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.Gateway.GRPCNetwork,
						"address": cfg.Reva.Gateway.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"gateway": map[string]interface{}{
								// registries is located on the gateway
								"authregistrysvc":    cfg.Reva.Gateway.Endpoint,
								"storageregistrysvc": cfg.Reva.Gateway.Endpoint,
								"appregistrysvc":     cfg.Reva.Gateway.Endpoint,
								// user metadata is located on the users services
								"preferencessvc":   cfg.Reva.Users.Endpoint,
								"userprovidersvc":  cfg.Reva.Users.Endpoint,
								"groupprovidersvc": cfg.Reva.Groups.Endpoint,
								// sharing is located on the sharing service
								"usershareprovidersvc":          cfg.Reva.Sharing.Endpoint,
								"publicshareprovidersvc":        cfg.Reva.Sharing.Endpoint,
								"ocmshareprovidersvc":           cfg.Reva.Sharing.Endpoint,
								"commit_share_to_storage_grant": cfg.Reva.Gateway.CommitShareToStorageGrant,
								"commit_share_to_storage_ref":   cfg.Reva.Gateway.CommitShareToStorageRef,
								"share_folder":                  cfg.Reva.Gateway.ShareFolder, // ShareFolder is the location where to create shares in the recipient's storage provider.
								// other
								"disable_home_creation_on_login": cfg.Reva.Gateway.DisableHomeCreationOnLogin,
								"datagateway":                    cfg.Reva.DataGateway.PublicURL,
								"transfer_shared_secret":         cfg.Reva.TransferSecret,
								"transfer_expires":               cfg.Reva.TransferExpires,
								"home_mapping":                   cfg.Reva.Gateway.HomeMapping,
								"etag_cache_ttl":                 cfg.Reva.Gateway.EtagCacheTTL,
							},
							"authregistry": map[string]interface{}{
								"driver": "static",
								"drivers": map[string]interface{}{
									"static": map[string]interface{}{
										"rules": map[string]interface{}{
											"basic":        cfg.Reva.AuthBasic.Endpoint,
											"bearer":       cfg.Reva.AuthBearer.Endpoint,
											"publicshares": cfg.Reva.StoragePublicLink.Endpoint,
										},
									},
								},
							},
							"storageregistry": map[string]interface{}{
								"driver": cfg.Reva.StorageRegistry.Driver,
								"drivers": map[string]interface{}{
									"static": map[string]interface{}{
										"home_provider": cfg.Reva.StorageRegistry.HomeProvider,
										"rules":         rules(cfg),
									},
								},
							},
						},
					},
				}

				gr.Add(func() error {
					err := external.RegisterGRPCEndpoint(
						ctx,
						"com.owncloud.storage",
						uuid.String(),
						cfg.Reva.Gateway.GRPCAddr,
						logger,
					)

					if err != nil {
						return err
					}

					runtime.RunWithOptions(
						rcfg,
						pidFile,
						runtime.WithLogger(&logger.Logger),
					)
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

func rules(cfg *config.Config) map[string]interface{} {

	// if a list of rules is given it overrides the generated rules from below
	if len(cfg.Reva.StorageRegistry.Rules) > 0 {
		rules := map[string]interface{}{}
		for i := range cfg.Reva.StorageRegistry.Rules {
			parts := strings.SplitN(cfg.Reva.StorageRegistry.Rules[i], "=", 2)
			rules[parts[0]] = parts[1]
		}
		return rules
	}

	// generate rules based on default config
	return map[string]interface{}{
		cfg.Reva.StorageHome.MountPath:       cfg.Reva.StorageHome.Endpoint,
		cfg.Reva.StorageHome.MountID:         cfg.Reva.StorageHome.Endpoint,
		cfg.Reva.StorageUsers.MountPath:      cfg.Reva.StorageUsers.Endpoint,
		cfg.Reva.StorageUsers.MountID:        cfg.Reva.StorageUsers.Endpoint,
		cfg.Reva.StoragePublicLink.MountPath: cfg.Reva.StoragePublicLink.Endpoint,
		// public link storage returns the mount id of the actual storage
		// medatada storage not part of the global namespace
	}
}

// GatewaySutureService allows for the storage-gateway command to be embedded and supervised by a suture supervisor tree.
type GatewaySutureService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewGatewaySutureService creates a new gateway.GatewaySutureService
func NewGateway(ctx context.Context) GatewaySutureService {
	sctx, cancel := context.WithCancel(ctx)
	cfg := config.New()
	cfg.Context = sctx
	return GatewaySutureService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg,
	}
}

func (s GatewaySutureService) Serve() {
	f := &flag.FlagSet{}
	for k := range Gateway(s.cfg).Flags {
		if err := Gateway(s.cfg).Flags[k].Apply(f); err != nil {
			return
		}
	}
	ctx := cli.NewContext(nil, f, nil)
	if Gateway(s.cfg).Before != nil {
		if err := Gateway(s.cfg).Before(ctx); err != nil {
			return
		}
	}
	if err := Gateway(s.cfg).Action(ctx); err != nil {
		return
	}
}

func (s GatewaySutureService) Stop() {
	s.cancel()
}
