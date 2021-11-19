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
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// StoragePublicLink is the entrypoint for the reva-storage-public-link command.
func StoragePublicLink(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "storage-public-link",
		Usage: "Start storage-public-link service",
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg, "storage-public-link")
		},
		Category: "Extensions",
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.Must(uuid.NewV4()).String()+".pid")
			rcfg := storagePublicLinkConfigFromStruct(c, cfg)

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
				debug.Addr(cfg.Reva.StoragePublicLink.DebugAddr),
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

			if !cfg.Reva.StoragePublicLink.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// storagePublicLinkConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func storagePublicLinkConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"max_cpus":             cfg.Reva.StoragePublicLink.MaxCPUs,
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
			"network": cfg.Reva.StoragePublicLink.GRPCNetwork,
			"address": cfg.Reva.StoragePublicLink.GRPCAddr,
			"interceptors": map[string]interface{}{
				"log": map[string]interface{}{},
			},
			"services": map[string]interface{}{
				"publicstorageprovider": map[string]interface{}{
					"mount_path":   cfg.Reva.StoragePublicLink.MountPath,
					"mount_id":     cfg.Reva.StoragePublicLink.MountID,
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
	return rcfg
}

// StoragePublicLinkSutureService allows for the storage-public-link command to be embedded and supervised by a suture supervisor tree.
type StoragePublicLinkSutureService struct {
	cfg *config.Config
}

// NewStoragePublicLinkSutureService creates a new storage.StoragePublicLinkSutureService
func NewStoragePublicLink(cfg *ociscfg.Config) suture.Service {
	cfg.Storage.Commons = cfg.Commons
	return StoragePublicLinkSutureService{
		cfg: cfg.Storage,
	}
}

func (s StoragePublicLinkSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.StoragePublicLink.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := StoragePublicLink(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if StoragePublicLink(s.cfg).Before != nil {
		if err := StoragePublicLink(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := StoragePublicLink(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
