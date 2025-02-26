package command

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/config"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/logging"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/revaconfig"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/server/http"
	"github.com/urfave/cli/v2"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(_ *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			if cfg.AllowImpersonation {
				fmt.Println("WARNING: Impersonation is enabled. Admins can impersonate all users.")
			}

			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}
			gr := run.Group{}
			ctx, cancel := context.WithCancel(c.Context)

			defer cancel()

			// make sure the run group executes all interrupt handlers when the context is canceled
			gr.Add(func() error {
				<-ctx.Done()
				return nil
			}, func(_ error) {
			})

			gr.Add(func() error {
				pidFile := path.Join(os.TempDir(), "revad-"+cfg.Service.Name+"-"+uuid.Must(uuid.NewV4()).String()+".pid")
				rCfg := revaconfig.AuthAppConfigFromStruct(cfg)
				reg := registry.GetRegistry()

				runtime.RunWithOptions(rCfg, pidFile,
					runtime.WithLogger(&logger.Logger),
					runtime.WithRegistry(reg),
					runtime.WithTraceProvider(traceProvider),
				)

				return nil
			}, func(err error) {
				if err == nil {
					logger.Info().
						Str("transport", "reva").
						Str("server", cfg.Service.Name).
						Msg("Shutting down server")
				} else {
					logger.Error().Err(err).
						Str("transport", "reva").
						Str("server", cfg.Service.Name).
						Msg("Shutting down server")
				}

				cancel()
			})

			debugServer, err := debug.Server(
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)
			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				_ = debugServer.Shutdown(ctx)
				cancel()
			})

			grpcSvc := registry.BuildGRPCService(cfg.GRPC.Namespace+"."+cfg.Service.Name, cfg.GRPC.Protocol, cfg.GRPC.Addr, version.GetString())
			if err := registry.RegisterService(ctx, logger, grpcSvc, cfg.Debug.Addr); err != nil {
				logger.Fatal().Err(err).Msg("failed to register the grpc service")
			}

			tm, err := pool.StringToTLSMode(cfg.GRPCClientTLS.Mode)
			if err != nil {
				return err
			}
			gatewaySelector, err := pool.GatewaySelector(
				cfg.Reva.Address,
				append(
					cfg.Reva.GetRevaOptions(),
					pool.WithTLSCACert(cfg.GRPCClientTLS.CACert),
					pool.WithTLSMode(tm),
					pool.WithRegistry(registry.GetRegistry()),
					pool.WithTracerProvider(traceProvider),
				)...)
			if err != nil {
				return err
			}

			grpcClient, err := ogrpc.NewClient(
				append(ogrpc.GetClientOptions(cfg.GRPCClientTLS), ogrpc.WithTraceProvider(traceProvider))...,
			)
			if err != nil {
				return err
			}

			rClient := settingssvc.NewRoleService("com.owncloud.api.settings", grpcClient)
			server, err := http.Server(
				http.Logger(logger),
				http.Context(ctx),
				http.Config(cfg),
				http.GatewaySelector(gatewaySelector),
				http.RoleClient(rClient),
				http.TracerProvider(traceProvider),
			)
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to initialize http server")
			}

			gr.Add(server.Run, func(err error) {
				logger.Error().Err(err).Str("server", "http").Msg("shutting down server")
			})

			return gr.Run()
		},
	}
}
