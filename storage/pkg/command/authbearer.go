package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// AuthBearer is the entrypoint for the auth-bearer command.
func AuthBearer(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-bearer",
		Usage: "Start authprovider for bearer auth",
		Flags: flagset.AuthBearerWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.AuthBearer.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)
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
				debug.Addr(cfg.Reva.AuthBearer.DebugAddr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Reva.AuthBearer.Supervised {
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
			"max_cpus":             cfg.Reva.AuthBearer.MaxCPUs,
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
			"network": cfg.Reva.AuthBearer.GRPCNetwork,
			"address": cfg.Reva.AuthBearer.GRPCAddr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"authprovider": map[string]interface{}{
					"auth_manager": "oidc",
					"auth_managers": map[string]interface{}{
						"oidc": map[string]interface{}{
							"issuer":     cfg.Reva.OIDC.Issuer,
							"insecure":   cfg.Reva.OIDC.Insecure,
							"id_claim":   cfg.Reva.OIDC.IDClaim,
							"uid_claim":  cfg.Reva.OIDC.UIDClaim,
							"gid_claim":  cfg.Reva.OIDC.GIDClaim,
							"gatewaysvc": cfg.Reva.Gateway.Endpoint,
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
	cfg.Storage.Log = cfg.Commons.Log
	return AuthBearerSutureService{
		cfg: cfg.Storage,
	}
}

func (s AuthBearerSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.AuthBearer.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := AuthBearer(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if AuthBearer(s.cfg).Before != nil {
		if err := AuthBearer(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := AuthBearer(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
