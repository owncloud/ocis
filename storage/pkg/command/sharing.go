package command

import (
	"context"
	"flag"
	"os"
	"path"
	"path/filepath"

	"github.com/owncloud/ocis/storage/pkg/tracing"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Sharing is the entrypoint for the sharing command.
func Sharing(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "sharing",
		Usage: "Start sharing service",
		Flags: flagset.SharingWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Sharing.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			tracing.Configure(cfg, logger)

			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()

			// precreate folders
			if cfg.Reva.Sharing.UserDriver == "json" && cfg.Reva.Sharing.UserJSONFile != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.Sharing.UserJSONFile), os.FileMode(0700)); err != nil {
					return err
				}
			}
			if cfg.Reva.Sharing.PublicDriver == "json" && cfg.Reva.Sharing.PublicJSONFile != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.Sharing.PublicJSONFile), os.FileMode(0700)); err != nil {
					return err
				}
			}

			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

			rcfg := sharingConfigFromStruct(c, cfg)

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

			debug, err := debug.Server(
				debug.Name(c.Command.Name+"-debug"),
				debug.Addr(cfg.Reva.Sharing.DebugAddr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", c.Command.Name+"-debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debug.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Reva.Sharing.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// sharingConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func sharingConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"max_cpus":             cfg.Reva.Sharing.MaxCPUs,
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
			"network": cfg.Reva.Sharing.GRPCNetwork,
			"address": cfg.Reva.Sharing.GRPCAddr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"usershareprovider": map[string]interface{}{
					"driver": cfg.Reva.Sharing.UserDriver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file": cfg.Reva.Sharing.UserJSONFile,
						},
						"sql": map[string]interface{}{ // cernbox sql
							"db_username":                   cfg.Reva.Sharing.UserSQLUsername,
							"db_password":                   cfg.Reva.Sharing.UserSQLPassword,
							"db_host":                       cfg.Reva.Sharing.UserSQLHost,
							"db_port":                       cfg.Reva.Sharing.UserSQLPort,
							"db_name":                       cfg.Reva.Sharing.UserSQLName,
							"password_hash_cost":            cfg.Reva.Sharing.PublicPasswordHashCost,
							"enable_expired_shares_cleanup": cfg.Reva.Sharing.PublicEnableExpiredSharesCleanup,
							"janitor_run_interval":          cfg.Reva.Sharing.PublicJanitorRunInterval,
						},
						"oc10-sql": map[string]interface{}{
							"gateway_addr":     cfg.Reva.Gateway.Endpoint,
							"storage_mount_id": cfg.Reva.Sharing.UserStorageMountID,
							"db_username":      cfg.Reva.Sharing.UserSQLUsername,
							"db_password":      cfg.Reva.Sharing.UserSQLPassword,
							"db_host":          cfg.Reva.Sharing.UserSQLHost,
							"db_port":          cfg.Reva.Sharing.UserSQLPort,
							"db_name":          cfg.Reva.Sharing.UserSQLName,
						},
					},
				},
				"publicshareprovider": map[string]interface{}{
					"driver": cfg.Reva.Sharing.PublicDriver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file": cfg.Reva.Sharing.PublicJSONFile,
						},
						"sql": map[string]interface{}{
							"db_username":                   cfg.Reva.Sharing.UserSQLUsername,
							"db_password":                   cfg.Reva.Sharing.UserSQLPassword,
							"db_host":                       cfg.Reva.Sharing.UserSQLHost,
							"db_port":                       cfg.Reva.Sharing.UserSQLPort,
							"db_name":                       cfg.Reva.Sharing.UserSQLName,
							"password_hash_cost":            cfg.Reva.Sharing.PublicPasswordHashCost,
							"enable_expired_shares_cleanup": cfg.Reva.Sharing.PublicEnableExpiredSharesCleanup,
							"janitor_run_interval":          cfg.Reva.Sharing.PublicJanitorRunInterval,
						},
						"oc10-sql": map[string]interface{}{
							"gateway_addr":                  cfg.Reva.Gateway.Endpoint,
							"storage_mount_id":              cfg.Reva.Sharing.UserStorageMountID,
							"db_username":                   cfg.Reva.Sharing.UserSQLUsername,
							"db_password":                   cfg.Reva.Sharing.UserSQLPassword,
							"db_host":                       cfg.Reva.Sharing.UserSQLHost,
							"db_port":                       cfg.Reva.Sharing.UserSQLPort,
							"db_name":                       cfg.Reva.Sharing.UserSQLName,
							"password_hash_cost":            cfg.Reva.Sharing.PublicPasswordHashCost,
							"enable_expired_shares_cleanup": cfg.Reva.Sharing.PublicEnableExpiredSharesCleanup,
							"janitor_run_interval":          cfg.Reva.Sharing.PublicJanitorRunInterval,
						},
					},
				},
			},
		},
	}
	return rcfg
}

// SharingSutureService allows for the storage-sharing command to be embedded and supervised by a suture supervisor tree.
type SharingSutureService struct {
	cfg *config.Config
}

// NewSharingSutureService creates a new store.SharingSutureService
func NewSharing(cfg *ociscfg.Config) suture.Service {
	if cfg.Mode == 0 {
		cfg.Storage.Reva.Sharing.Supervised = true
	}
	return SharingSutureService{
		cfg: cfg.Storage,
	}
}

func (s SharingSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.Sharing.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := Sharing(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if Sharing(s.cfg).Before != nil {
		if err := Sharing(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := Sharing(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
