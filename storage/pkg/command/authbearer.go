package command

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
)

// AuthBearer is the entrypoint for the auth-bearer command.
func AuthBearer(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-bearer",
		Usage: "Start authprovider for bearer auth",
		Flags: flagset.AuthBearerWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.AuthBearer.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					logger.Error().
						Str("type", t).
						Msg("Storage only supports the jaeger tracing backend")

				case "jaeger":
					logger.Info().
						Str("type", t).
						Msg("configuring storage to use the jaeger tracing backend")

				case "zipkin":
					logger.Error().
						Str("type", t).
						Msg("Storage only supports the jaeger tracing backend")

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
						"max_cpus":             cfg.Reva.AuthBearer.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": c.Command.Name,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.AuthBearer.GRPCNetwork,
						"address": cfg.Reva.AuthBearer.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"authprovider": map[string]interface{}{
								"auth_manager": "oidc",
								"auth_managers": map[string]interface{}{
									"oidc": map[string]interface{}{
										"issuer":     cfg.Reva.OIDC.Issuer,
										"insecure":   cfg.Reva.OIDC.Insecure,
										"id_claim":   cfg.Reva.OIDC.IDClaim,
										"uid_claim":  cfg.Reva.OIDC.UIDClaim,
										"gid_claim":  cfg.Reva.OIDC.GIDClaim,
										"gatewaysvc": cfg.Reva.Gateway.Endpoint,
									},
								},
							},
						},
					},
				}

				gr.Add(func() error {
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
					debug.Addr(cfg.Reva.AuthBearer.DebugAddr),
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

// AuthBearerSutureService allows for the storage-gateway command to be embedded and supervised by a suture supervisor tree.
type AuthBearerSutureService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewAuthBearerSutureService creates a new gateway.AuthBearerSutureService
func NewAuthBearer(ctx context.Context, cfg *config.Config) AuthBearerSutureService {
	sctx, cancel := context.WithCancel(ctx)
	cfg.Context = sctx
	return AuthBearerSutureService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg,
	}
}

func (s AuthBearerSutureService) Serve() {
	f := &flag.FlagSet{}
	for k := range AuthBearer(s.cfg).Flags {
		if err := AuthBearer(s.cfg).Flags[k].Apply(f); err != nil {
			return
		}
	}
	ctx := cli.NewContext(nil, f, nil)
	if AuthBearer(s.cfg).Before != nil {
		if err := AuthBearer(s.cfg).Before(ctx); err != nil {
			return
		}
	}
	if err := AuthBearer(s.cfg).Action(ctx); err != nil {
		return
	}
}

func (s AuthBearerSutureService) Stop() {
	s.cancel()
}
