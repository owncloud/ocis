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
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/storage/pkg/command/storagedrivers"
	"github.com/owncloud/ocis/storage/pkg/config"
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
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg, "storage-metadata")
		},
		Category: "Extensions",
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
				version.String,
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
			"jwt_secret":                cfg.Reva.JWTSecret,
			"gatewaysvc":                cfg.Reva.Gateway.Endpoint,
			"skip_user_groups_in_token": cfg.Reva.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.Reva.StorageMetadata.GRPCNetwork,
			"address": cfg.Reva.StorageMetadata.GRPCAddr,
			"interceptors": map[string]interface{}{
				"log": map[string]interface{}{},
			},
			"services": map[string]interface{}{
				"storageprovider": map[string]interface{}{
					"driver":          cfg.Reva.StorageMetadata.Driver,
					"drivers":         storagedrivers.MetadataDrivers(cfg),
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
					"drivers":     storagedrivers.MetadataDrivers(cfg),
					"timeout":     86400,
					"insecure":    cfg.Reva.StorageMetadata.DataProvider.Insecure,
					"disable_tus": true,
				},
			},
		},
	}
	return rcfg
}

// SutureService allows for the storage-metadata command to be embedded and supervised by a suture supervisor tree.
type MetadataSutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new storagemetadata.SutureService
func NewStorageMetadata(cfg *ociscfg.Config) suture.Service {
	cfg.Storage.Commons = cfg.Commons
	return MetadataSutureService{
		cfg: cfg.Storage,
	}
}

func (s MetadataSutureService) Serve(ctx context.Context) error {
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
