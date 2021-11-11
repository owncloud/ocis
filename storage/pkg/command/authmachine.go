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

// AuthMachine is the entrypoint for the auth-machine command.
func AuthMachine(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-machine",
		Usage: "Start authprovider for machine auth",
		Flags: flagset.AuthMachineWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.AuthMachine.Services = c.StringSlice("service")

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
			rcfg := authMachineConfigFromStruct(c, cfg)

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
				debug.Addr(cfg.Reva.AuthMachine.DebugAddr),
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

			if !cfg.Reva.AuthMachine.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// authMachineConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func authMachineConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"core": map[string]interface{}{
			"max_cpus":             cfg.Reva.AuthMachine.MaxCPUs,
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
			"network": cfg.Reva.AuthMachine.GRPCNetwork,
			"address": cfg.Reva.AuthMachine.GRPCAddr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"authprovider": map[string]interface{}{
					"auth_manager": "machine",
					"auth_managers": map[string]interface{}{
						"machine": map[string]interface{}{
							"api_key":      cfg.Reva.AuthMachineConfig.MachineAuthAPIKey,
							"gateway_addr": cfg.Reva.Gateway.Endpoint,
						},
					},
				},
			},
		},
	}
}

// AuthMachineSutureService allows for the storage-gateway command to be embedded and supervised by a suture supervisor tree.
type AuthMachineSutureService struct {
	cfg *config.Config
}

// NewAuthMachineSutureService creates a new gateway.AuthMachineSutureService
func NewAuthMachine(cfg *ociscfg.Config) suture.Service {
	cfg.Storage.Log = cfg.Commons.Log
	return AuthMachineSutureService{
		cfg: cfg.Storage,
	}
}

func (s AuthMachineSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.AuthMachine.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := AuthMachine(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if AuthMachine(s.cfg).Before != nil {
		if err := AuthMachine(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := AuthMachine(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
