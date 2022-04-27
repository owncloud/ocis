package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/auth-bearer/pkg/config"
	"github.com/owncloud/ocis/extensions/auth-bearer/pkg/config/parser"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// AuthBearer is the entrypoint for the auth-bearer command.
func AuthBearer(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-bearer",
		Usage: "start authprovider for bearer auth",
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

			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")
			rcfg := authBearerConfigFromStruct(c, cfg)

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
				logger.Info().Err(err).Str("server", "debug").Msg("failed to initialize server")
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

// authBearerConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func authBearerConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
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
				"authprovider": map[string]interface{}{
					"auth_manager": cfg.AuthProvider,
					"auth_managers": map[string]interface{}{
						"oidc": map[string]interface{}{
							"issuer":    cfg.AuthProviders.OIDC.Issuer,
							"insecure":  cfg.AuthProviders.OIDC.Insecure,
							"id_claim":  cfg.AuthProviders.OIDC.IDClaim,
							"uid_claim": cfg.AuthProviders.OIDC.UIDClaim,
							"gid_claim": cfg.AuthProviders.OIDC.GIDClaim,
						},
					},
				},
			},
		},
	}
}

// AuthBearerSutureService allows for the storage-gateway command to be embedded and supervised by a suture supervisor tree.
type AuthBearerSutureService struct {
	cfg *config.Config
}

// NewAuthBearerSutureService creates a new gateway.AuthBearerSutureService
func NewAuthBearer(cfg *ociscfg.Config) suture.Service {
	cfg.AuthBearer.Commons = cfg.Commons
	return AuthBearerSutureService{
		cfg: cfg.AuthBearer,
	}
}

func (s AuthBearerSutureService) Serve(ctx context.Context) error {
	cmd := AuthBearer(s.cfg)
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
