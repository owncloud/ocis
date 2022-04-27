package command

import (
	"context"
	"flag"
	"os"
	"path"
	"path/filepath"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/sharing/pkg/config"
	"github.com/owncloud/ocis/extensions/sharing/pkg/config/parser"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Sharing is the entrypoint for the sharing command.
func Sharing(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "sharing",
		Usage: "start sharing service",
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
			if cfg.UserSharingDriver == "json" && cfg.UserSharingDrivers.JSON.File != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.UserSharingDrivers.JSON.File), os.FileMode(0700)); err != nil {
					return err
				}
			}
			if cfg.PublicSharingDriver == "json" && cfg.PublicSharingDrivers.JSON.File != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.PublicSharingDrivers.JSON.File), os.FileMode(0700)); err != nil {
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

			gr.Add(debug.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Supervised {
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
				"usershareprovider": map[string]interface{}{
					"driver": cfg.UserSharingDriver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file":         cfg.UserSharingDrivers.JSON.File,
							"gateway_addr": cfg.Reva.Address,
						},
						"sql": map[string]interface{}{ // cernbox sql
							"db_username":                   cfg.UserSharingDrivers.SQL.DBUsername,
							"db_password":                   cfg.UserSharingDrivers.SQL.DBPassword,
							"db_host":                       cfg.UserSharingDrivers.SQL.DBHost,
							"db_port":                       cfg.UserSharingDrivers.SQL.DBPort,
							"db_name":                       cfg.UserSharingDrivers.SQL.DBName,
							"password_hash_cost":            cfg.UserSharingDrivers.SQL.PasswordHashCost,
							"enable_expired_shares_cleanup": cfg.UserSharingDrivers.SQL.EnableExpiredSharesCleanup,
							"janitor_run_interval":          cfg.UserSharingDrivers.SQL.JanitorRunInterval,
						},
						"oc10-sql": map[string]interface{}{
							"storage_mount_id": cfg.UserSharingDrivers.SQL.UserStorageMountID,
							"db_username":      cfg.UserSharingDrivers.SQL.DBUsername,
							"db_password":      cfg.UserSharingDrivers.SQL.DBPassword,
							"db_host":          cfg.UserSharingDrivers.SQL.DBHost,
							"db_port":          cfg.UserSharingDrivers.SQL.DBPort,
							"db_name":          cfg.UserSharingDrivers.SQL.DBName,
						},
						"cs3": map[string]interface{}{
							"provider_addr":       cfg.UserSharingDrivers.CS3.ProviderAddr,
							"service_user_id":     cfg.UserSharingDrivers.CS3.ServiceUserID,
							"service_user_idp":    cfg.UserSharingDrivers.CS3.ServiceUserIDP,
							"machine_auth_apikey": cfg.UserSharingDrivers.CS3.MachineAuthAPIKey,
						},
					},
				},
				"publicshareprovider": map[string]interface{}{
					"driver": cfg.PublicSharingDriver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file":         cfg.PublicSharingDrivers.JSON.File,
							"gateway_addr": cfg.Reva.Address,
						},
						"sql": map[string]interface{}{
							"db_username":                   cfg.PublicSharingDrivers.SQL.DBUsername,
							"db_password":                   cfg.PublicSharingDrivers.SQL.DBPassword,
							"db_host":                       cfg.PublicSharingDrivers.SQL.DBHost,
							"db_port":                       cfg.PublicSharingDrivers.SQL.DBPort,
							"db_name":                       cfg.PublicSharingDrivers.SQL.DBName,
							"password_hash_cost":            cfg.PublicSharingDrivers.SQL.PasswordHashCost,
							"enable_expired_shares_cleanup": cfg.PublicSharingDrivers.SQL.EnableExpiredSharesCleanup,
							"janitor_run_interval":          cfg.PublicSharingDrivers.SQL.JanitorRunInterval,
						},
						"oc10-sql": map[string]interface{}{
							"storage_mount_id":              cfg.PublicSharingDrivers.SQL.UserStorageMountID,
							"db_username":                   cfg.PublicSharingDrivers.SQL.DBUsername,
							"db_password":                   cfg.PublicSharingDrivers.SQL.DBPassword,
							"db_host":                       cfg.PublicSharingDrivers.SQL.DBHost,
							"db_port":                       cfg.PublicSharingDrivers.SQL.DBPort,
							"db_name":                       cfg.PublicSharingDrivers.SQL.DBName,
							"password_hash_cost":            cfg.PublicSharingDrivers.SQL.PasswordHashCost,
							"enable_expired_shares_cleanup": cfg.PublicSharingDrivers.SQL.EnableExpiredSharesCleanup,
							"janitor_run_interval":          cfg.PublicSharingDrivers.SQL.JanitorRunInterval,
						},
						"cs3": map[string]interface{}{
							"provider_addr":       cfg.PublicSharingDrivers.CS3.ProviderAddr,
							"service_user_id":     cfg.PublicSharingDrivers.CS3.ServiceUserID,
							"service_user_idp":    cfg.PublicSharingDrivers.CS3.ServiceUserIDP,
							"machine_auth_apikey": cfg.PublicSharingDrivers.CS3.MachineAuthAPIKey,
						},
					},
				},
			},
			"interceptors": map[string]interface{}{
				"eventsmiddleware": map[string]interface{}{
					"group":     "sharing",
					"type":      "nats",
					"address":   cfg.Events.Addr,
					"clusterID": cfg.Events.ClusterID,
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
	cfg.Sharing.Commons = cfg.Commons
	return SharingSutureService{
		cfg: cfg.Sharing,
	}
}

func (s SharingSutureService) Serve(ctx context.Context) error {
	// s.cfg.Reva.Sharing.Context = ctx
	cmd := Sharing(s.cfg)
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
