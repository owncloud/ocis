package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// StorageShares is the entrypoint for the reva-storage-share command.
func StorageShares(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-shares",
		Usage:    "Start storage-shares service",
		Flags:    flagset.StorageShares(cfg),
		Category: "Extensions",
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.Must(uuid.NewV4()).String()+".pid")
			rcfg := storageSharesConfigFromStruct(c, cfg)

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
				debug.Addr(cfg.Reva.StorageShares.DebugAddr),
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

			if !cfg.Reva.StorageMetadata.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// storageSharesConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func storageSharesConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"max_cpus":             cfg.Reva.StorageShares.MaxCPUs,
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
			"network": cfg.Reva.StorageShares.GRPCNetwork,
			"address": cfg.Reva.StorageShares.GRPCAddr,
			"interceptors": map[string]interface{}{
				"log": map[string]interface{}{},
			},
			"services": map[string]interface{}{
				"sharesstorageprovider": map[string]interface{}{
					"mount_path":           cfg.Reva.StorageShares.MountPath,
					"gateway_addr":         cfg.Reva.Gateway.Endpoint,
					"usershareprovidersvc": cfg.Reva.Sharing.Endpoint,
				},
			},
		},
	}
	return rcfg
}

// StorageSharesSutureService allows for the storage-shares command to be embedded and supervised by a suture supervisor tree.
type StorageSharesSutureService struct {
	cfg *config.Config
}

// NewStorageSharesSutureService creates a new storage.StorageSharesSutureService
func NewStorageShares(cfg *ociscfg.Config) suture.Service {
	if cfg.Mode == 0 {
		cfg.Storage.Reva.StorageShares.Supervised = true
	}
	return StorageSharesSutureService{
		cfg: cfg.Storage,
	}
}

func (s StorageSharesSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.StorageShares.Context = ctx
	f := &flag.FlagSet{}
	for k := range StorageShares(s.cfg).Flags {
		if err := StorageShares(s.cfg).Flags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if StorageShares(s.cfg).Before != nil {
		if err := StorageShares(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := StorageShares(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
