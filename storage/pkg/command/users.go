package command

import (
	"context"
	"flag"
	"os"
	"path"
	"path/filepath"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
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

// Users is the entrypoint for the users command.
func Users(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "start users service",
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg, "storage-users")
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			tracing.Configure(cfg, logger)

			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()

			// precreate folders
			if cfg.Reva.Users.Driver == "json" && cfg.Reva.Users.JSON != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.Users.JSON), os.FileMode(0700)); err != nil {
					return err
				}
			}

			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

			rcfg := usersConfigFromStruct(c, cfg)

			logger.Debug().
				Str("server", "users").
				Interface("reva-config", rcfg).
				Msg("config")

			if cfg.Reva.Users.Driver == "ldap" {
				if err := waitForLDAPCA(logger, &cfg.Reva.LDAP); err != nil {
					logger.Error().Err(err).Msg("The configured LDAP CA cert does not exist")
					return err
				}
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

			debugServer, err := debug.Server(
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

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Reva.Users.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// usersConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func usersConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"max_cpus":             cfg.Reva.Users.MaxCPUs,
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
						"ldap": ldapConfigFromString(cfg),
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
						"owncloudsql": map[string]interface{}{
							"dbusername":           cfg.Reva.UserOwnCloudSQL.DBUsername,
							"dbpassword":           cfg.Reva.UserOwnCloudSQL.DBPassword,
							"dbhost":               cfg.Reva.UserOwnCloudSQL.DBHost,
							"dbport":               cfg.Reva.UserOwnCloudSQL.DBPort,
							"dbname":               cfg.Reva.UserOwnCloudSQL.DBName,
							"idp":                  cfg.Reva.UserOwnCloudSQL.Idp,
							"nobody":               cfg.Reva.UserOwnCloudSQL.Nobody,
							"join_username":        cfg.Reva.UserOwnCloudSQL.JoinUsername,
							"join_ownclouduuid":    cfg.Reva.UserOwnCloudSQL.JoinOwnCloudUUID,
							"enable_medial_search": cfg.Reva.UserOwnCloudSQL.EnableMedialSearch,
						},
					},
				},
			},
		},
	}
	return rcfg
}

// UserProviderSutureService allows for the storage-userprovider command to be embedded and supervised by a suture supervisor tree.
type UserProviderSutureService struct {
	cfg *config.Config
}

// NewUserProviderSutureService creates a new storage.UserProvider
func NewUserProvider(cfg *ociscfg.Config) suture.Service {
	cfg.Storage.Commons = cfg.Commons
	return UserProviderSutureService{
		cfg: cfg.Storage,
	}
}

func (s UserProviderSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.Users.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := Users(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if Users(s.cfg).Before != nil {
		if err := Users(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := Users(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
