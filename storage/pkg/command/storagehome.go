package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/tracing"
	"github.com/thejerf/suture/v4"
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

			tracing.Configure(cfg, logger)

			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

			// override driver enable home option with home config
			if cfg.Reva.Storages.Home.EnableHome {
				cfg.Reva.Storages.Common.EnableHome = true
				cfg.Reva.Storages.EOS.EnableHome = true
				cfg.Reva.Storages.Local.EnableHome = true
				cfg.Reva.Storages.OwnCloud.EnableHome = true
				cfg.Reva.Storages.S3.EnableHome = true
				cfg.Reva.Storages.S3NG.EnableHome = true
			}
			rcfg := storageHomeConfigFromStruct(c, cfg)

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

			debugServer, err := debug.Server(
				debug.Name(c.Command.Name+"-debug"),
				debug.Addr(cfg.Reva.StorageHome.DebugAddr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", c.Command.Name+"-debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Reva.StorageHome.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// storageHomeConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func storageHomeConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
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
	return rcfg
}

// StorageHomeSutureService allows for the storage-home command to be embedded and supervised by a suture supervisor tree.
type StorageHomeSutureService struct {
	cfg *config.Config
}

// NewStorageHomeSutureService creates a new storage.StorageHomeSutureService
func NewStorageHome(cfg *ociscfg.Config) suture.Service {
	if cfg.Mode == 0 {
		cfg.Storage.Reva.StorageHome.Supervised = true
	}
	return StorageHomeSutureService{
		cfg: cfg.Storage,
	}
}

func (s StorageHomeSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.StorageHome.Context = ctx
	f := &flag.FlagSet{}
	for k := range StorageHome(s.cfg).Flags {
		if err := StorageHome(s.cfg).Flags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if StorageHome(s.cfg).Before != nil {
		if err := StorageHome(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := StorageHome(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
