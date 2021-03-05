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

// StorageHome is the entrypoint for the storage-home command.
func StorageHome(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "storage-home",
		Usage: "Start storage-home service",
		Flags: flagset.StorageHomeWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.StorageHome.Services = c.StringSlice("service")

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

				// override driver enable home option with home config
				if cfg.Reva.Storages.Home.EnableHome {
					cfg.Reva.Storages.Common.EnableHome = true
					cfg.Reva.Storages.EOS.EnableHome = true
					cfg.Reva.Storages.Local.EnableHome = true
					cfg.Reva.Storages.OwnCloud.EnableHome = true
					cfg.Reva.Storages.S3.EnableHome = true
				}
				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.StorageHome.MaxCPUs,
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
						"network": cfg.Reva.StorageHome.GRPCNetwork,
						"address": cfg.Reva.StorageHome.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"storageprovider": map[string]interface{}{
								"driver":             cfg.Reva.StorageHome.Driver,
								"drivers":            drivers(cfg),
								"mount_path":         cfg.Reva.StorageHome.MountPath,
								"mount_id":           cfg.Reva.StorageHome.MountID,
								"expose_data_server": cfg.Reva.StorageHome.ExposeDataServer,
								"data_server_url":    cfg.Reva.StorageHome.DataServerURL,
								"tmp_folder":         cfg.Reva.StorageHome.TempFolder,
							},
						},
					},
					"http": map[string]interface{}{
						"network": cfg.Reva.StorageHome.HTTPNetwork,
						"address": cfg.Reva.StorageHome.HTTPAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"dataprovider": map[string]interface{}{
								"prefix":      cfg.Reva.StorageHome.HTTPPrefix,
								"driver":      cfg.Reva.StorageHome.Driver,
								"drivers":     drivers(cfg),
								"timeout":     86400,
								"insecure":    true,
								"disable_tus": false,
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
					debug.Addr(cfg.Reva.StorageHome.DebugAddr),
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

// StorageHomeSutureService allows for the storage-home command to be embedded and supervised by a suture supervisor tree.
type StorageHomeSutureService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewStorageHomeSutureService creates a new storage.StorageHomeSutureService
func NewStorageHome(ctx context.Context) StorageHomeSutureService {
	sctx, cancel := context.WithCancel(ctx)
	cfg := config.New()
	cfg.Context = sctx
	return StorageHomeSutureService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg,
	}
}

func (s StorageHomeSutureService) Serve() {
	f := &flag.FlagSet{}
	for k := range StorageHome(s.cfg).Flags {
		if err := StorageHome(s.cfg).Flags[k].Apply(f); err != nil {
			return
		}
	}
	ctx := cli.NewContext(nil, f, nil)
	if StorageHome(s.cfg).Before != nil {
		if err := StorageHome(s.cfg).Before(ctx); err != nil {
			return
		}
	}
	if err := StorageHome(s.cfg).Action(ctx); err != nil {
		return
	}
}

func (s StorageHomeSutureService) Stop() {
	s.cancel()
}
