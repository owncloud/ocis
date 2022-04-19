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
	"github.com/owncloud/ocis/extensions/auth-basic/pkg/config"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/ldap"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Command is the entrypoint for the auth-basic command.
func AuthBasic(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-basic",
		Usage: "start authprovider for basic auth",
		// Before: func(c *cli.Context) error {
		// 	return ParseConfig(c, cfg, "storage-auth-basic")
		// },
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

			// pre-create folders
			if cfg.AuthProvider == "json" && cfg.AuthProviders.JSON.File != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.AuthProviders.JSON.File), os.FileMode(0700)); err != nil {
					return err
				}
			}

			uuid := uuid.Must(uuid.NewV4())

			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

			rcfg := authBasicConfigFromStruct(c, cfg)
			logger.Debug().
				Str("server", "authbasic").
				Interface("reva-config", rcfg).
				Msg("config")

			if cfg.AuthProvider == "ldap" {
				ldapCfg := cfg.AuthProviders.LDAP
				if err := ldap.WaitForCA(logger, ldapCfg.Insecure, ldapCfg.CACert); err != nil {
					logger.Error().Err(err).Msg("The configured LDAP CA cert does not exist")
					return err
				}
			}

			gr.Add(func() error {
				runtime.RunWithOptions(rcfg, pidFile, runtime.WithLogger(&logger.Logger))
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
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
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

// authBasicConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func authBasicConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.JWTSecret,
			"gatewaysvc":                cfg.GatewayEndpoint,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"authprovider": map[string]interface{}{
					"auth_manager": cfg.AuthProvider,
					"auth_managers": map[string]interface{}{
						"json": map[string]interface{}{
							"users": cfg.AuthProviders.JSON.File,
						},
						"ldap": ldapConfigFromString(cfg.AuthProviders.LDAP),
						"owncloudsql": map[string]interface{}{
							"dbusername":        cfg.AuthProviders.OwnCloudSQL.DBUsername,
							"dbpassword":        cfg.AuthProviders.OwnCloudSQL.DBPassword,
							"dbhost":            cfg.AuthProviders.OwnCloudSQL.DBHost,
							"dbport":            cfg.AuthProviders.OwnCloudSQL.DBPort,
							"dbname":            cfg.AuthProviders.OwnCloudSQL.DBName,
							"idp":               cfg.AuthProviders.OwnCloudSQL.IDP,
							"nobody":            cfg.AuthProviders.OwnCloudSQL.Nobody,
							"join_username":     cfg.AuthProviders.OwnCloudSQL.JoinUsername,
							"join_ownclouduuid": cfg.AuthProviders.OwnCloudSQL.JoinOwnCloudUUID,
						},
					},
				},
			},
		},
	}
	return rcfg
}

// AuthBasicSutureService allows for the storage-authbasic command to be embedded and supervised by a suture supervisor tree.
type AuthBasicSutureService struct {
	cfg *config.Config
}

// NewAuthBasicSutureService creates a new store.AuthBasicSutureService
func NewAuthBasic(cfg *ociscfg.Config) suture.Service {
	cfg.AuthBasic.Commons = cfg.Commons
	return AuthBasicSutureService{
		cfg: cfg.AuthBasic,
	}
}

func (s AuthBasicSutureService) Serve(ctx context.Context) error {
	// s.cfg.Reva.AuthBasic.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := AuthBasic(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if AuthBasic(s.cfg).Before != nil {
		if err := AuthBasic(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := AuthBasic(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}

func ldapConfigFromString(cfg config.LDAPProvider) map[string]interface{} {
	return map[string]interface{}{
		"uri":               cfg.URI,
		"cacert":            cfg.CACert,
		"insecure":          cfg.Insecure,
		"bind_username":     cfg.BindDN,
		"bind_password":     cfg.BindPassword,
		"user_base_dn":      cfg.UserBaseDN,
		"group_base_dn":     cfg.GroupBaseDN,
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
