package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/storage-metadata/pkg/config"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	"github.com/owncloud/ocis/extensions/storage/pkg/service/external"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// StorageMetadata the entrypoint for the storage-storage-metadata command.
//
// It provides a ocis-specific storage store metadata (shares,account,settings...)
func StorageMetadata(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-metadata",
		Usage:    "start storage-metadata service",
		Category: "extensions",
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
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
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
				debug.Addr(cfg.Debug.Addr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Pprof(cfg.Debug.Pprof),
				debug.Zpages(cfg.Debug.Zpages),
				debug.Token(cfg.Debug.Token),
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

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			if err := external.RegisterGRPCEndpoint(
				ctx,
				"com.owncloud.storage.metadata",
				uuid.Must(uuid.NewV4()).String(),
				cfg.GRPC.Addr,
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
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.JWTSecret,
			"gatewaysvc":                cfg.GRPC.Addr,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"services": map[string]interface{}{
				"gateway": map[string]interface{}{
					// registries are located on the gateway
					"authregistrysvc":    cfg.GRPC.Addr,
					"storageregistrysvc": cfg.GRPC.Addr,
					// user metadata is located on the users services
					"userprovidersvc":  cfg.GRPC.Addr,
					"groupprovidersvc": cfg.GRPC.Addr,
					"permissionssvc":   cfg.GRPC.Addr,
					// other
					"disable_home_creation_on_login": true, // metadata manually creates a space
					// metadata always uses the simple upload, so no transfer secret or datagateway needed
				},
				"userprovider": map[string]interface{}{
					"driver": "memory",
					"drivers": map[string]interface{}{
						"memory": map[string]interface{}{
							"users": map[string]interface{}{
								"serviceuser": map[string]interface{}{
									"id": map[string]interface{}{
										"opaqueId": "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
										"idp":      "internal",
										"type":     1, // user.UserType_USER_TYPE_PRIMARY
									},
									"username":     "serviceuser",
									"display_name": "System User",
								},
							},
						},
					},
				},
				"authregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"rules": map[string]interface{}{
								"machine": cfg.GRPC.Addr,
							},
						},
					},
				},
				"authprovider": map[string]interface{}{
					"auth_manager": "machine",
					"auth_managers": map[string]interface{}{
						"machine": map[string]interface{}{
							"api_key":      cfg.MachineAuthAPIKey,
							"gateway_addr": cfg.GRPC.Addr,
						},
					},
				},
				"permissions": map[string]interface{}{
					"driver": "demo",
					"drivers": map[string]interface{}{
						"demo": map[string]interface{}{},
					},
				},
				"storageregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"rules": map[string]interface{}{
								"/": map[string]interface{}{
									"address": cfg.GRPC.Addr,
								},
							},
						},
					},
				},
				"storageprovider": map[string]interface{}{
					"driver":          cfg.Driver,
					"drivers":         config.MetadataDrivers(cfg),
					"data_server_url": cfg.DataServerURL,
					"tmp_folder":      cfg.TempFolder,
				},
			},
		},
		"http": map[string]interface{}{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			// no datagateway needed as the metadata clients directly talk to the dataprovider with the simple protocol
			"services": map[string]interface{}{
				"dataprovider": map[string]interface{}{
					"prefix":      "data",
					"driver":      cfg.Driver,
					"drivers":     config.MetadataDrivers(cfg),
					"timeout":     86400,
					"insecure":    cfg.DataProviderInsecure,
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
	cfg.StorageMetadata.Commons = cfg.Commons
	return MetadataSutureService{
		cfg: cfg.StorageMetadata,
	}
}

func (s MetadataSutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
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
