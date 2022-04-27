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
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	"github.com/owncloud/ocis/extensions/user/pkg/config"
	"github.com/owncloud/ocis/extensions/user/pkg/config/parser"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/ldap"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// User is the entrypoint for the user command.
func User(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "start users service",
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

			// precreate folders
			if cfg.Driver == "json" && cfg.Drivers.JSON.File != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Drivers.JSON.File), os.FileMode(0700)); err != nil {
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

			if cfg.Driver == "ldap" {
				if err := ldap.WaitForCA(logger, cfg.Drivers.LDAP.Insecure, cfg.Drivers.LDAP.CACert); err != nil {
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

// usersConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func usersConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
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
			// TODO build services dynamically
			"services": map[string]interface{}{
				"userprovider": map[string]interface{}{
					"driver": cfg.Driver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"users": cfg.Drivers.JSON.File,
						},
						"ldap": ldapConfigFromString(cfg.Drivers.LDAP),
						"rest": map[string]interface{}{
							"client_id":                    cfg.Drivers.REST.ClientID,
							"client_secret":                cfg.Drivers.REST.ClientSecret,
							"redis_address":                cfg.Drivers.REST.RedisAddr,
							"redis_username":               cfg.Drivers.REST.RedisUsername,
							"redis_password":               cfg.Drivers.REST.RedisPassword,
							"user_groups_cache_expiration": cfg.UsersCacheExpiration,
							"id_provider":                  cfg.Drivers.REST.IDProvider,
							"api_base_url":                 cfg.Drivers.REST.APIBaseURL,
							"oidc_token_endpoint":          cfg.Drivers.REST.OIDCTokenEndpoint,
							"target_api":                   cfg.Drivers.REST.TargetAPI,
						},
						"owncloudsql": map[string]interface{}{
							"dbusername":           cfg.Drivers.OwnCloudSQL.DBUsername,
							"dbpassword":           cfg.Drivers.OwnCloudSQL.DBPassword,
							"dbhost":               cfg.Drivers.OwnCloudSQL.DBHost,
							"dbport":               cfg.Drivers.OwnCloudSQL.DBPort,
							"dbname":               cfg.Drivers.OwnCloudSQL.DBName,
							"idp":                  cfg.Drivers.OwnCloudSQL.IDP,
							"nobody":               cfg.Drivers.OwnCloudSQL.Nobody,
							"join_username":        cfg.Drivers.OwnCloudSQL.JoinUsername,
							"join_ownclouduuid":    cfg.Drivers.OwnCloudSQL.JoinOwnCloudUUID,
							"enable_medial_search": cfg.Drivers.OwnCloudSQL.EnableMedialSearch,
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
	cfg.User.Commons = cfg.Commons
	return UserProviderSutureService{
		cfg: cfg.User,
	}
}

func (s UserProviderSutureService) Serve(ctx context.Context) error {
	// s.cfg.Reva.Users.Context = ctx
	cmd := User(s.cfg)
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

func ldapConfigFromString(cfg config.LDAPDriver) map[string]interface{} {
	return map[string]interface{}{
		"uri":               cfg.URI,
		"cacert":            cfg.CACert,
		"insecure":          cfg.Insecure,
		"bind_username":     cfg.BindDN,
		"bind_password":     cfg.BindPassword,
		"user_base_dn":      cfg.UserBaseDN,
		"group_base_dn":     cfg.GroupBaseDN,
		"user_scope":        cfg.UserScope,
		"group_scope":       cfg.GroupScope,
		"user_filter":       cfg.UserFilter,
		"group_filter":      cfg.GroupFilter,
		"user_objectclass":  cfg.UserObjectClass,
		"group_objectclass": cfg.GroupObjectClass,
		"login_attributes":  cfg.LoginAttributes,
		"idp":               cfg.IDP,
		"user_schema": map[string]interface{}{
			"id":              cfg.UserSchema.ID,
			"idIsOctetString": cfg.UserSchema.IDIsOctetString,
			"mail":            cfg.UserSchema.Mail,
			"displayName":     cfg.UserSchema.DisplayName,
			"userName":        cfg.UserSchema.Username,
		},
		"group_schema": map[string]interface{}{
			"id":              cfg.GroupSchema.ID,
			"idIsOctetString": cfg.GroupSchema.IDIsOctetString,
			"mail":            cfg.GroupSchema.Mail,
			"displayName":     cfg.GroupSchema.DisplayName,
			"groupName":       cfg.GroupSchema.Groupname,
			"member":          cfg.GroupSchema.Member,
		},
	}
}
