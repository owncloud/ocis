package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/thejerf/suture"
)

// StorageUsers is the entrypoint for the storage-users command.
func StorageUsers(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "storage-users",
		Usage: "Start storage-users service",
		Flags: flagset.StorageUsersWithConfig(cfg),
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
						"max_cpus":             cfg.Reva.StorageUsers.MaxCPUs,
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
						"network": cfg.Reva.StorageUsers.GRPCNetwork,
						"address": cfg.Reva.StorageUsers.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"storageprovider": map[string]interface{}{
								"driver":             cfg.Reva.StorageUsers.Driver,
								"drivers":            drivers(cfg),
								"mount_path":         cfg.Reva.StorageUsers.MountPath,
								"mount_id":           cfg.Reva.StorageUsers.MountID,
								"expose_data_server": cfg.Reva.StorageUsers.ExposeDataServer,
								"data_server_url":    cfg.Reva.StorageUsers.DataServerURL,
								"tmp_folder":         cfg.Reva.StorageUsers.TempFolder,
							},
						},
					},
					"http": map[string]interface{}{
						"network": cfg.Reva.StorageUsers.HTTPNetwork,
						"address": cfg.Reva.StorageUsers.HTTPAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"dataprovider": map[string]interface{}{
								"prefix":      cfg.Reva.StorageUsers.HTTPPrefix,
								"driver":      cfg.Reva.StorageUsers.Driver,
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
					debug.Addr(cfg.Reva.StorageUsers.DebugAddr),
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().Err(err).Str("server", c.Command.Name+"-debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(server.ListenAndServe, func(_ error) {
					cancel()
				})
			}

			if !cfg.Reva.StorageMetadata.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// StorageUsersSutureService allows for the storage-home command to be embedded and supervised by a suture supervisor tree.
type StorageUsersSutureService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewStorageUsersSutureService creates a new storage.StorageUsersSutureService
func NewStorageUsers(ctx context.Context, cfg *ociscfg.Config) suture.Service {
	sctx, cancel := context.WithCancel(ctx)
	cfg.Storage.Reva.StorageUsers.Context = sctx
	if cfg.Mode == 0 {
		cfg.Storage.Reva.StorageUsers.Supervised = true
	}
	return StorageUsersSutureService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg.Storage,
	}
}

func (s StorageUsersSutureService) Serve() {
	f := &flag.FlagSet{}
	for k := range StorageUsers(s.cfg).Flags {
		if err := StorageUsers(s.cfg).Flags[k].Apply(f); err != nil {
			return
		}
	}
	ctx := cli.NewContext(nil, f, nil)
	if StorageUsers(s.cfg).Before != nil {
		if err := StorageUsers(s.cfg).Before(ctx); err != nil {
			return
		}
	}
	if err := StorageUsers(s.cfg).Action(ctx); err != nil {
		return
	}
}

func (s StorageUsersSutureService) Stop() {
	s.cancel()
}
