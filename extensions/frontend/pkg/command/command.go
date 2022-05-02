package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/frontend/pkg/config"
	"github.com/owncloud/ocis/extensions/frontend/pkg/config/parser"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Frontend is the entrypoint for the frontend command.
func Frontend(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "frontend",
		Usage: "start frontend service",
		Before: func(ctx *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
			}
			return err
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
			//metrics     = metrics.New()

			defer cancel()

			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

			archivers := []map[string]interface{}{
				{
					"enabled":       true,
					"version":       "2.0.0",
					"formats":       []string{"tar", "zip"},
					"archiver_url":  path.Join("/", cfg.Archiver.Prefix),
					"max_num_files": strconv.FormatInt(cfg.Archiver.MaxNumFiles, 10),
					"max_size":      strconv.FormatInt(cfg.Archiver.MaxSize, 10),
				},
			}

			appProviders := []map[string]interface{}{
				{
					"enabled":  true,
					"version":  "1.0.0",
					"apps_url": cfg.AppProvider.AppsURL,
					"open_url": cfg.AppProvider.OpenURL,
					"new_url":  cfg.AppProvider.NewURL,
				},
			}

			filesCfg := map[string]interface{}{
				"private_links":     false,
				"bigfilechunking":   false,
				"blacklisted_files": []string{},
				"undelete":          true,
				"versioning":        true,
				"archivers":         archivers,
				"app_providers":     appProviders,
				"favorites":         cfg.EnableFavorites,
			}

			if cfg.DefaultUploadProtocol == "tus" {
				filesCfg["tus_support"] = map[string]interface{}{
					"version":              "1.0.0",
					"resumable":            "1.0.0",
					"extension":            "creation,creation-with-upload",
					"http_method_override": cfg.UploadHTTPMethodOverride,
					"max_chunk_size":       cfg.UploadMaxChunkSize,
				}
			}

			revaCfg := frontendConfigFromStruct(c, cfg, filesCfg)

			gr.Add(func() error {
				runtime.RunWithOptions(revaCfg, pidFile, runtime.WithLogger(&logger.Logger))
				return nil
			}, func(_ error) {
				logger.Info().Str("server", c.Command.Name).Msg("Shutting down server")
				cancel()
			})

			{
				server, err := debug.Server(
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
						Str("server", "debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(server.ListenAndServe, func(_ error) {
					cancel()
				})
			}

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// frontendConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func frontendConfigFromStruct(c *cli.Context, cfg *config.Config, filesCfg map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address, // Todo or address?
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"http": map[string]interface{}{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			"middlewares": map[string]interface{}{
				"cors": map[string]interface{}{
					"allow_credentials": true,
				},
				"auth": map[string]interface{}{
					"credentials_by_user_agent": cfg.Middleware.Auth.CredentialsByUserAgent,
					"credential_chain":          []string{"bearer"},
				},
			},
			// TODO build services dynamically
			"services": map[string]interface{}{
				"appprovider": map[string]interface{}{
					"prefix":                 cfg.AppProvider.Prefix,
					"transfer_shared_secret": cfg.TransferSecret,
					"timeout":                86400,
					"insecure":               cfg.AppProvider.Insecure,
				},
				"archiver": map[string]interface{}{
					"prefix":        cfg.Archiver.Prefix,
					"timeout":       86400,
					"insecure":      cfg.Archiver.Insecure,
					"max_num_files": cfg.Archiver.MaxNumFiles,
					"max_size":      cfg.Archiver.MaxSize,
				},
				"datagateway": map[string]interface{}{
					"prefix":                 cfg.DataGateway.Prefix,
					"transfer_shared_secret": cfg.TransferSecret,
					"timeout":                86400,
					"insecure":               true,
				},
				"ocs": map[string]interface{}{
					"storage_registry_svc":      cfg.Reva.Address,
					"share_prefix":              cfg.OCS.SharePrefix,
					"home_namespace":            cfg.OCS.HomeNamespace,
					"resource_info_cache_ttl":   cfg.OCS.ResourceInfoCacheTTL,
					"prefix":                    cfg.OCS.Prefix,
					"additional_info_attribute": cfg.OCS.AdditionalInfoAttribute,
					"machine_auth_apikey":       cfg.MachineAuthAPIKey,
					"cache_warmup_driver":       cfg.OCS.CacheWarmupDriver,
					"cache_warmup_drivers": map[string]interface{}{
						"cbox": map[string]interface{}{
							"db_username": cfg.OCS.CacheWarmupDrivers.CBOX.DBUsername,
							"db_password": cfg.OCS.CacheWarmupDrivers.CBOX.DBPassword,
							"db_host":     cfg.OCS.CacheWarmupDrivers.CBOX.DBHost,
							"db_port":     cfg.OCS.CacheWarmupDrivers.CBOX.DBPort,
							"db_name":     cfg.OCS.CacheWarmupDrivers.CBOX.DBName,
							"namespace":   cfg.OCS.CacheWarmupDrivers.CBOX.Namespace,
							"gatewaysvc":  cfg.Reva.Address,
						},
					},
					"config": map[string]interface{}{
						"version": "1.7",
						"website": "ownCloud",
						"host":    cfg.PublicURL,
						"contact": "",
						"ssl":     "false",
					},
					"default_upload_protocol": cfg.DefaultUploadProtocol,
					"capabilities": map[string]interface{}{
						"capabilities": map[string]interface{}{
							"core": map[string]interface{}{
								"poll_interval": 60,
								"webdav_root":   "remote.php/webdav",
								"status": map[string]interface{}{
									"installed":      true,
									"maintenance":    false,
									"needsDbUpgrade": false,
									"version":        "10.0.11.5",
									"versionstring":  "10.0.11",
									"edition":        "community",
									"productname":    "reva",
									"hostname":       "",
								},
								"support_url_signing": true,
							},
							"checksums": map[string]interface{}{
								"supported_types":       cfg.Checksums.SupportedTypes,
								"preferred_upload_type": cfg.Checksums.PreferredUploadType,
							},
							"files": filesCfg,
							"dav": map[string]interface{}{
								"reports": []string{"search-files"},
							},
							"files_sharing": map[string]interface{}{
								"api_enabled":                       true,
								"resharing":                         false,
								"group_sharing":                     true,
								"auto_accept_share":                 true,
								"share_with_group_members_only":     true,
								"share_with_membership_groups_only": true,
								"default_permissions":               22,
								"search_min_length":                 3,
								"public": map[string]interface{}{
									"enabled":              true,
									"send_mail":            true,
									"social_share":         true,
									"upload":               true,
									"multiple":             true,
									"supports_upload_only": true,
									"password": map[string]interface{}{
										"enforced": false,
										"enforced_for": map[string]interface{}{
											"read_only":   false,
											"read_write":  false,
											"upload_only": false,
										},
									},
									"expire_date": map[string]interface{}{
										"enabled": true,
									},
									"can_edit": true,
								},
								"user": map[string]interface{}{
									"send_mail":       true,
									"profile_picture": false,
									"settings": []map[string]interface{}{
										{
											"enabled": true,
											"version": "1.0.0",
										},
									},
								},
								"user_enumeration": map[string]interface{}{
									"enabled":            true,
									"group_members_only": true,
								},
								"federation": map[string]interface{}{
									"outgoing": true,
									"incoming": true,
								},
							},
							"spaces": map[string]interface{}{
								"version":    "0.0.1",
								"enabled":    cfg.EnableProjectSpaces || cfg.EnableShareJail,
								"projects":   cfg.EnableProjectSpaces,
								"share_jail": cfg.EnableShareJail,
							},
						},
						"version": map[string]interface{}{
							"edition": "reva",
							"major":   10,
							"minor":   0,
							"micro":   11,
							"string":  "10.0.11",
						},
					},
				},
			},
		},
	}
}

// FrontendSutureService allows for the storage-frontend command to be embedded and supervised by a suture supervisor tree.
type FrontendSutureService struct {
	cfg *config.Config
}

// NewFrontend creates a new frontend.FrontendSutureService
func NewFrontend(cfg *ociscfg.Config) suture.Service {
	cfg.Frontend.Commons = cfg.Commons
	return FrontendSutureService{
		cfg: cfg.Frontend,
	}
}

func (s FrontendSutureService) Serve(ctx context.Context) error {
	// s.cfg.Reva.Frontend.Context = ctx
	cmd := Frontend(s.cfg)
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
