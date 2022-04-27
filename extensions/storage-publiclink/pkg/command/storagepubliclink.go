package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/storage-publiclink/pkg/config"
	"github.com/owncloud/ocis/extensions/storage-publiclink/pkg/config/parser"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// StoragePublicLink is the entrypoint for the reva-storage-public-link command.
func StoragePublicLink(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-public-link",
		Usage:    "start storage-public-link service",
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Action: func(c *cli.Context) error {
			logCfg := cfg.Logging
			logger := log.NewLogger(
				log.Level(logCfg.Level),
				log.File(logCfg.File),
				log.Pretty(logCfg.Pretty),
				log.Color(logCfg.Color),
			)
			tracing.Configure(cfg.Tracing.Enabled, cfg.Tracing.Type, logger)
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
				debug.Addr(cfg.Debug.Addr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Pprof(cfg.Debug.Pprof),
				debug.Zpages(cfg.Debug.Zpages),
				debug.Token(cfg.Debug.Token),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", c.Command.Name+"-debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Supervised {
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
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"interceptors": map[string]interface{}{
				"log": map[string]interface{}{},
			},
			"services": map[string]interface{}{
				"publicstorageprovider": map[string]interface{}{
					"mount_id":     cfg.StorageProvider.MountID,
					"gateway_addr": cfg.StorageProvider.GatewayEndpoint,
				},
				"authprovider": map[string]interface{}{
					"auth_manager": "publicshares",
					"auth_managers": map[string]interface{}{
						"publicshares": map[string]interface{}{
							"gateway_addr": cfg.AuthProvider.GatewayEndpoint,
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
	cfg.StoragePublicLink.Commons = cfg.Commons
	return StoragePublicLinkSutureService{
		cfg: cfg.StoragePublicLink,
	}
}

func (s StoragePublicLinkSutureService) Serve(ctx context.Context) error {
	// s.cfg.Reva.StoragePublicLink.Context = ctx
	cmd := StoragePublicLink(s.cfg)
	f := &flag.FlagSet{}
	cmdFlags := cmd.Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if cmd.Before != nil {
		if err := cmd.Before(cliCtx); err != nil {
			return err
		}
	}
	if err := cmd.Action(cliCtx); err != nil {
		return err
	}

	return nil
}
