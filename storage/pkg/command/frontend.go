package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/conversions"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Frontend is the entrypoint for the frontend command.
func Frontend(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "frontend",
		Usage: "Start frontend service",
		Before: func(c *cli.Context) error {
			if err := loadUserAgent(c, cfg); err != nil {
				return err
			}
			return ParseConfig(c, cfg, "storage-frontend")
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			tracing.Configure(cfg, logger)

			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			//metrics     = metrics.New()

			defer cancel()

			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

			// pregenerate list of valid localhost ports for the desktop redirect_uri
			// TODO use custom scheme like "owncloud://localhost/user/callback" tracked in
			var desktopRedirectURIs [65535 - 1024]string
			for port := 0; port < len(desktopRedirectURIs); port++ {
				desktopRedirectURIs[port] = fmt.Sprintf("http://localhost:%d", (port + 1024))
			}

			archivers := []map[string]interface{}{
				{
					"enabled":       true,
					"version":       "2.0.0",
					"formats":       []string{"tar", "zip"},
					"archiver_url":  cfg.Reva.Archiver.ArchiverURL,
					"max_num_files": strconv.FormatInt(cfg.Reva.Archiver.MaxNumFiles, 10),
					"max_size":      strconv.FormatInt(cfg.Reva.Archiver.MaxSize, 10),
				},
			}

			appProviders := []map[string]interface{}{
				{
					"enabled":  true,
					"version":  "1.0.0",
					"apps_url": cfg.Reva.AppProvider.AppsURL,
					"open_url": cfg.Reva.AppProvider.OpenURL,
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
				"favorites":         cfg.Reva.Frontend.Favorites,
			}

			if cfg.Reva.DefaultUploadProtocol == "tus" {
				filesCfg["tus_support"] = map[string]interface{}{
					"version":              "1.0.0",
					"resumable":            "1.0.0",
					"extension":            "creation,creation-with-upload",
					"http_method_override": cfg.Reva.UploadHTTPMethodOverride,
					"max_chunk_size":       cfg.Reva.UploadMaxChunkSize,
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
					debug.Addr(cfg.Reva.Frontend.DebugAddr),
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
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

			if !cfg.Reva.Frontend.Supervised {
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
			"max_cpus":             cfg.Reva.Users.MaxCPUs,
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.Reva.JWTSecret,
			"gatewaysvc":                cfg.Reva.Gateway.Endpoint, // Todo or address?
			"skip_user_groups_in_token": cfg.Reva.SkipUserGroupsInToken,
		},
		"http": map[string]interface{}{
			"network": cfg.Reva.Frontend.HTTPNetwork,
			"address": cfg.Reva.Frontend.HTTPAddr,
			"middlewares": map[string]interface{}{
				"cors": map[string]interface{}{
					"allow_credentials": true,
				},
				"auth": map[string]interface{}{
					"credentials_by_user_agent": cfg.Reva.Frontend.Middleware.Auth.CredentialsByUserAgent,
					"credential_chain":          []string{"bearer"},
				},
			},
			// TODO build services dynamically
			"services": map[string]interface{}{
				"appprovider": map[string]interface{}{
					"prefix":                 cfg.Reva.Frontend.AppProviderPrefix,
					"transfer_shared_secret": cfg.Reva.TransferSecret,
					"timeout":                86400,
					"insecure":               cfg.Reva.Frontend.AppProviderInsecure,
				},
				"archiver": map[string]interface{}{
					"prefix":        cfg.Reva.Frontend.ArchiverPrefix,
					"timeout":       86400,
					"insecure":      cfg.Reva.Frontend.ArchiverInsecure,
					"max_num_files": cfg.Reva.Archiver.MaxNumFiles,
					"max_size":      cfg.Reva.Archiver.MaxSize,
				},
				"datagateway": map[string]interface{}{
					"prefix":                 cfg.Reva.Frontend.DatagatewayPrefix,
					"transfer_shared_secret": cfg.Reva.TransferSecret,
					"timeout":                86400,
					"insecure":               true,
				},
				"ocdav": map[string]interface{}{
					"prefix":           cfg.Reva.Frontend.OCDavPrefix,
					"files_namespace":  cfg.Reva.OCDav.DavFilesNamespace,
					"webdav_namespace": cfg.Reva.OCDav.WebdavNamespace,
					"timeout":          86400,
					"insecure":         cfg.Reva.Frontend.OCDavInsecure,
					"public_url":       cfg.Reva.Frontend.PublicURL,
				},
				"ocs": map[string]interface{}{
					"share_prefix":              cfg.Reva.Frontend.OCSSharePrefix,
					"home_namespace":            cfg.Reva.Frontend.OCSHomeNamespace,
					"resource_info_cache_ttl":   cfg.Reva.Frontend.OCSResourceInfoCacheTTL,
					"prefix":                    cfg.Reva.Frontend.OCSPrefix,
					"additional_info_attribute": cfg.Reva.Frontend.OCSAdditionalInfoAttribute,
					"cache_warmup_driver":       cfg.Reva.Frontend.OCSCacheWarmupDriver,
					"cache_warmup_drivers": map[string]interface{}{
						"cbox": map[string]interface{}{
							"db_username": cfg.Reva.Sharing.UserSQLUsername,
							"db_password": cfg.Reva.Sharing.UserSQLPassword,
							"db_host":     cfg.Reva.Sharing.UserSQLHost,
							"db_port":     cfg.Reva.Sharing.UserSQLPort,
							"db_name":     cfg.Reva.Sharing.UserSQLName,
							"namespace":   cfg.Reva.UserStorage.EOS.Root,
							"gatewaysvc":  cfg.Reva.Gateway.Endpoint,
						},
					},
					"config": map[string]interface{}{
						"version": "1.8",
						"website": "reva",
						"host":    cfg.Reva.Frontend.PublicURL,
						"contact": "admin@localhost",
						"ssl":     "false",
					},
					"default_upload_protocol": cfg.Reva.DefaultUploadProtocol,
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
								"supported_types":       cfg.Reva.ChecksumSupportedTypes,
								"preferred_upload_type": cfg.Reva.ChecksumPreferredUploadType,
							},
							"files": filesCfg,
							"dav":   map[string]interface{}{},
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
										"enforced": true,
										"enforced_for": map[string]interface{}{
											"read_only":   true,
											"read_write":  true,
											"upload_only": true,
										},
									},
									"expire_date": map[string]interface{}{
										"enabled": false,
									},
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

// loadUserAgent reads the user-agent-whitelist-lock-in, since it is a string flag, and attempts to construct a map of
// "user-agent":"challenge" locks in for Reva.
// Modifies cfg. Spaces don't need to be trimmed as urfavecli takes care of it. User agents with spaces are valid. i.e:
// Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:83.0) Gecko/20100101 Firefox/83.0
// This function works by relying in our format of specifying [user-agent:challenge] and the fact that the user agent
// might contain ":" (colon), so the original string is reversed, split in two parts, by the time it is split we
// have the indexes reversed and the tuple is in the format of [challenge:user-agent], then the same process is applied
// in reverse for each individual part
func loadUserAgent(c *cli.Context, cfg *config.Config) error {
	cfg.Reva.Frontend.Middleware.Auth.CredentialsByUserAgent = make(map[string]string)
	locks := c.StringSlice("user-agent-whitelist-lock-in")

	for _, v := range locks {
		vv := conversions.Reverse(v)
		parts := strings.SplitN(vv, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("unexpected config value for user-agent lock-in: %v, expected format is user-agent:challenge", v)
		}

		cfg.Reva.Frontend.Middleware.Auth.CredentialsByUserAgent[conversions.Reverse(parts[1])] = conversions.Reverse(parts[0])
	}

	return nil
}

// FrontendSutureService allows for the storage-frontend command to be embedded and supervised by a suture supervisor tree.
type FrontendSutureService struct {
	cfg *config.Config
}

// NewFrontend creates a new frontend.FrontendSutureService
func NewFrontend(cfg *ociscfg.Config) suture.Service {
	cfg.Storage.Commons = cfg.Commons
	return FrontendSutureService{
		cfg: cfg.Storage,
	}
}

func (s FrontendSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.Frontend.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := Frontend(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if Frontend(s.cfg).Before != nil {
		if err := Frontend(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := Frontend(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
