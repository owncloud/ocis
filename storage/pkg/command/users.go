package command

import (
	"context"
	"flag"
	"os"
	"path"
	"path/filepath"

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

// Users is the entrypoint for the sharing command.
func Users(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Start users service",
		Flags: flagset.UsersWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Users.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					logger.Error().
						Str("type", t).
						Msg("Reva only supports the jaeger tracing backend")

				case "jaeger":
					logger.Info().
						Str("type", t).
						Msg("configuring storage to use the jaeger tracing backend")

				case "zipkin":
					logger.Error().
						Str("type", t).
						Msg("Reva only supports the jaeger tracing backend")

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

			// precreate folders
			if cfg.Reva.Users.Driver == "json" && cfg.Reva.Users.JSON != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.Users.JSON), os.FileMode(0700)); err != nil {
					return err
				}
			}

			{
				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.Users.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": c.Command.Name,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.Users.GRPCNetwork,
						"address": cfg.Reva.Users.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"userprovider": map[string]interface{}{
								"driver": cfg.Reva.Users.Driver,
								"drivers": map[string]interface{}{
									"json": map[string]interface{}{
										"users": cfg.Reva.Users.JSON,
									},
									"ldap": map[string]interface{}{
										"hostname":        cfg.Reva.LDAP.Hostname,
										"port":            cfg.Reva.LDAP.Port,
										"base_dn":         cfg.Reva.LDAP.BaseDN,
										"userfilter":      cfg.Reva.LDAP.UserFilter,
										"attributefilter": cfg.Reva.LDAP.UserAttributeFilter,
										"findfilter":      cfg.Reva.LDAP.UserFindFilter,
										"groupfilter":     cfg.Reva.LDAP.UserGroupFilter,
										"bind_username":   cfg.Reva.LDAP.BindDN,
										"bind_password":   cfg.Reva.LDAP.BindPassword,
										"idp":             cfg.Reva.LDAP.IDP,
										"schema": map[string]interface{}{
											"dn":          "dn",
											"uid":         cfg.Reva.LDAP.UserSchema.UID,
											"mail":        cfg.Reva.LDAP.UserSchema.Mail,
											"displayName": cfg.Reva.LDAP.UserSchema.DisplayName,
											"cn":          cfg.Reva.LDAP.UserSchema.CN,
											"uidNumber":   cfg.Reva.LDAP.UserSchema.UIDNumber,
											"gidNumber":   cfg.Reva.LDAP.UserSchema.GIDNumber,
										},
									},
									"rest": map[string]interface{}{
										"client_id":                    cfg.Reva.UserGroupRest.ClientID,
										"client_secret":                cfg.Reva.UserGroupRest.ClientSecret,
										"redis_address":                cfg.Reva.UserGroupRest.RedisAddress,
										"redis_username":               cfg.Reva.UserGroupRest.RedisUsername,
										"redis_password":               cfg.Reva.UserGroupRest.RedisPassword,
										"user_groups_cache_expiration": cfg.Reva.Users.UserGroupsCacheExpiration,
										"id_provider":                  cfg.Reva.UserGroupRest.IDProvider,
										"api_base_url":                 cfg.Reva.UserGroupRest.APIBaseURL,
										"oidc_token_endpoint":          cfg.Reva.UserGroupRest.OIDCTokenEndpoint,
										"target_api":                   cfg.Reva.UserGroupRest.TargetAPI,
									},
								},
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
					debug.Addr(cfg.Reva.Users.DebugAddr),
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

// UsersProviderService allows for the storage-userprovider command to be embedded and supervised by a suture supervisor tree.
type UsersProviderService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewUsersProviderService creates a new storage.UsersProviderService
func NewUsersProviderService(ctx context.Context, cfg *ociscfg.Config) suture.Service {
	sctx, cancel := context.WithCancel(ctx)
	cfg.Storage.Reva.Users.Context = sctx
	if cfg.Mode == 0 {
		cfg.Storage.Reva.Users.Supervised = true
	}
	return UsersProviderService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg.Storage,
	}
}

func (s UsersProviderService) Serve() {
	f := &flag.FlagSet{}
	for k := range Users(s.cfg).Flags {
		if err := Users(s.cfg).Flags[k].Apply(f); err != nil {
			return
		}
	}
	ctx := cli.NewContext(nil, f, nil)
	if Users(s.cfg).Before != nil {
		if err := Users(s.cfg).Before(ctx); err != nil {
			return
		}
	}
	if err := Users(s.cfg).Action(ctx); err != nil {
		return
	}
}

func (s UsersProviderService) Stop() {
	s.cancel()
}
