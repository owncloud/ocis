package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/service/external"
	"github.com/owncloud/ocis/storage/pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// StorageMetadata the entrypoint for the storage-storage-metadata command.
//
// It provides a ocis-specific storage store metadata (shares,account,settings...)
func StorageMetadata(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "storage-metadata",
		Usage: "Start storage-metadata service",
		// TODO(refs) at this point it might make sense delegate log flags to each individual storage command.
		Flags:    append(flagset.StorageMetadata(cfg), flagset.RootWithConfig(cfg)...),
		Category: "Extensions",
		Before: func(c *cli.Context) error {
			storageRoot := c.String("storage-root")

			cfg.Reva.Storages.OwnCloud.Root = storageRoot
			cfg.Reva.Storages.EOS.Root = storageRoot
			cfg.Reva.Storages.Local.Root = storageRoot
			cfg.Reva.Storages.S3.Root = storageRoot
			cfg.Reva.Storages.Home.Root = storageRoot

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)

			gr := run.Group{}
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Reva.StorageMetadata.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Reva.StorageMetadata.Context)
			}()

			defer cancel()

			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.Must(uuid.NewV4()).String()+".pid")

			// Disable home because the metadata is stored independently
			// of the user. This also means that a valid-token without any user-id
			// is allowed to write to the metadata-storage.
			cfg.Reva.Storages.Common.EnableHome = false
			cfg.Reva.Storages.EOS.EnableHome = false
			cfg.Reva.Storages.Local.EnableHome = false
			cfg.Reva.Storages.OwnCloud.EnableHome = false
			cfg.Reva.Storages.S3.EnableHome = false

			// We need this hack because the metadata storage can define STORAGE_METADATA_ROOT which has the same destination as
			// STORAGE_DRIVER_OCIS_ROOT. When both variables are set one storage will always be out of sync. Ensure the
			// metadata storage root is never overridden. This is the kind of stateful code that make you want to cry blood.
			if os.Getenv("STORAGE_METADATA_ROOT") != "" && os.Getenv("STORAGE_DRIVER_OCIS_ROOT") != "" {
				cfg.Reva.Storages.Common.Root = os.Getenv("STORAGE_METADATA_ROOT")
			}

			rcfg := storageMetadataFromStruct(c, cfg)

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
				debug.Addr(cfg.Reva.StorageMetadata.DebugAddr),
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
				return debugServer.ListenAndServe()
			}, func(_ error) {
				_ = debugServer.Shutdown(ctx)
				cancel()
			})

			if !cfg.Reva.StorageMetadata.Supervised {
				sync.Trap(&gr, cancel)
			}

			if err := external.RegisterGRPCEndpoint(
				ctx,
				"com.owncloud.storage.metadata",
				uuid.Must(uuid.NewV4()).String(),
				cfg.Reva.StorageMetadata.GRPCAddr,
				logger,
			); err != nil {
				logger.Fatal().Err(err).Msg("failed to register the grpc endpoint")
			}

			return gr.Run()
		},
	}
}

// storageMetadataFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func storageMetadataFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"max_cpus":             cfg.Reva.StorageMetadata.MaxCPUs,
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
			"network": cfg.Reva.StorageMetadata.GRPCNetwork,
			"address": cfg.Reva.StorageMetadata.GRPCAddr,
			"interceptors": map[string]interface{}{
				"log": map[string]interface{}{},
			},
			"services": map[string]interface{}{
				"storageprovider": map[string]interface{}{
					"mount_path":      "/meta",
					"driver":          cfg.Reva.StorageMetadata.Driver,
					"drivers":         drivers(cfg),
					"data_server_url": cfg.Reva.StorageMetadata.DataServerURL,
					"tmp_folder":      cfg.Reva.StorageMetadata.TempFolder,
				},
			},
		},
		"http": map[string]interface{}{
			"network": cfg.Reva.StorageMetadata.HTTPNetwork,
			"address": cfg.Reva.StorageMetadata.HTTPAddr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"dataprovider": map[string]interface{}{
					"prefix":      "data",
					"driver":      cfg.Reva.StorageMetadata.Driver,
					"drivers":     drivers(cfg),
					"timeout":     86400,
					"insecure":    true,
					"disable_tus": true,
				},
			},
		},
	}
	return rcfg
}

// SutureService allows for the storage-metadata command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new storagemetadata.SutureService
func NewStorageMetadata(cfg *ociscfg.Config) suture.Service {
	if cfg.Mode == 0 {
		cfg.Storage.Reva.StorageMetadata.Supervised = true
	}
	return SutureService{
		cfg: cfg.Storage,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.StorageMetadata.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := StorageMetadata(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if StorageMetadata(s.cfg).Before != nil {
		if err := StorageMetadata(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := StorageMetadata(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
