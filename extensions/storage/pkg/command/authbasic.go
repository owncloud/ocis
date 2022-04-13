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
	"github.com/owncloud/ocis/extensions/storage/pkg/config"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	"github.com/owncloud/ocis/extensions/storage/pkg/tracing"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// AuthBasic is the entrypoint for the auth-basic command.
func AuthBasic(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-basic",
		Usage: "start authprovider for basic auth",
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg, "storage-auth-basic")
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// pre-create folders
			if cfg.Reva.AuthProvider.Driver == "json" && cfg.Reva.AuthProvider.JSON != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.AuthProvider.JSON), os.FileMode(0700)); err != nil {
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

			if cfg.Reva.AuthProvider.Driver == "ldap" {
				if err := waitForLDAPCA(logger, &cfg.Reva.LDAP); err != nil {
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
				debug.Addr(cfg.Reva.AuthBasic.DebugAddr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Reva.AuthBasic.Supervised {
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
			"max_cpus":             cfg.Reva.AuthBasic.MaxCPUs,
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
			"network": cfg.Reva.AuthBasic.GRPCNetwork,
			"address": cfg.Reva.AuthBasic.GRPCAddr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"authprovider": map[string]interface{}{
					"auth_manager": cfg.Reva.AuthProvider.Driver,
					"auth_managers": map[string]interface{}{
						"json": map[string]interface{}{
							"users": cfg.Reva.AuthProvider.JSON,
						},
						"ldap": ldapConfigFromString(cfg),
						"owncloudsql": map[string]interface{}{
							"dbusername":        cfg.Reva.UserOwnCloudSQL.DBUsername,
							"dbpassword":        cfg.Reva.UserOwnCloudSQL.DBPassword,
							"dbhost":            cfg.Reva.UserOwnCloudSQL.DBHost,
							"dbport":            cfg.Reva.UserOwnCloudSQL.DBPort,
							"dbname":            cfg.Reva.UserOwnCloudSQL.DBName,
							"idp":               cfg.Reva.UserOwnCloudSQL.Idp,
							"nobody":            cfg.Reva.UserOwnCloudSQL.Nobody,
							"join_username":     cfg.Reva.UserOwnCloudSQL.JoinUsername,
							"join_ownclouduuid": cfg.Reva.UserOwnCloudSQL.JoinOwnCloudUUID,
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
	cfg.Storage.Commons = cfg.Commons
	return AuthBasicSutureService{
		cfg: cfg.Storage,
	}
}

func (s AuthBasicSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.AuthBasic.Context = ctx
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
