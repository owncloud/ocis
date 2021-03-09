package command

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path"
	"time"

	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/thejerf/suture"
	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
)

// StoragePublicLink is the entrypoint for the reva-storage-public-link command.
func StoragePublicLink(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-public-link",
		Usage:    "Start storage-public-link service",
		Flags:    flagset.StoragePublicLink(cfg),
		Category: "Extensions",
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
						"max_cpus":             cfg.Reva.StoragePublicLink.MaxCPUs,
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
						"network": cfg.Reva.StoragePublicLink.GRPCNetwork,
						"address": cfg.Reva.StoragePublicLink.GRPCAddr,
						"interceptors": map[string]interface{}{
							"log": map[string]interface{}{},
						},
						"services": map[string]interface{}{
							"publicstorageprovider": map[string]interface{}{
								"mount_path":   cfg.Reva.StoragePublicLink.MountPath,
								"gateway_addr": cfg.Reva.Gateway.Endpoint,
							},
							"authprovider": map[string]interface{}{
								"auth_manager": "publicshares",
								"auth_managers": map[string]interface{}{
									"publicshares": map[string]interface{}{
										"gateway_addr": cfg.Reva.Gateway.Endpoint,
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
					debug.Addr(cfg.Reva.StoragePublicLink.DebugAddr),
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("server", c.Command.Name+"-debug").
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
							Str("server", c.Command.Name+"-debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("server", c.Command.Name+"-debug").
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

// StoragePublicLinkSutureService allows for the storage-public-link command to be embedded and supervised by a suture supervisor tree.
type StoragePublicLinkSutureService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewStoragePublicLinkSutureService creates a new storage.StoragePublicLinkSutureService
func NewStoragePublicLink(ctx context.Context, cfg *ociscfg.Config) suture.Service {
	sctx, cancel := context.WithCancel(ctx)
	cfg.Storage.Reva.StoragePublicLink.Context = sctx
	return StoragePublicLinkSutureService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg.Storage,
	}
}

func (s StoragePublicLinkSutureService) Serve() {
	f := &flag.FlagSet{}
	for k := range StoragePublicLink(s.cfg).Flags {
		if err := StoragePublicLink(s.cfg).Flags[k].Apply(f); err != nil {
			return
		}
	}
	ctx := cli.NewContext(nil, f, nil)
	if StoragePublicLink(s.cfg).Before != nil {
		if err := StoragePublicLink(s.cfg).Before(ctx); err != nil {
			return
		}
	}
	if err := StoragePublicLink(s.cfg).Action(ctx); err != nil {
		return
	}
}

func (s StoragePublicLinkSutureService) Stop() {
	s.cancel()
}
