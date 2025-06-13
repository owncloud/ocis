package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/settings/pkg/logging"
	"github.com/owncloud/ocis/v2/services/settings/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/settings/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/settings/pkg/server/grpc"
	"github.com/owncloud/ocis/v2/services/settings/pkg/server/http"
	svc "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}
			cfg.GrpcClient, err = ogrpc.NewClient(
				append(ogrpc.GetClientOptions(cfg.GRPCClientTLS), ogrpc.WithTraceProvider(traceProvider))...,
			)
			if err != nil {
				return err
			}

			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			mtrcs := metrics.New()
			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			handle := svc.NewDefaultLanguageService(cfg, svc.NewService(cfg, logger))

			gr := runner.NewGroup()

			// prepare an HTTP server and add it to the group run.
			httpServer, err := http.Server(
				http.Name(cfg.Service.Name),
				http.Logger(logger),
				http.Context(ctx),
				http.Config(cfg),
				http.Metrics(mtrcs),
				http.ServiceHandler(handle),
				http.TraceProvider(traceProvider),
			)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("Error initializing http service")
				return fmt.Errorf("could not initialize http service: %w", err)
			}
			gr.Add(runner.NewGoMicroHttpServerRunner(cfg.Service.Name+".http", httpServer))

			// prepare a gRPC server and add it to the group run.
			grpcServer := grpc.Server(
				grpc.Name(cfg.Service.Name),
				grpc.Logger(logger),
				grpc.Context(ctx),
				grpc.Config(cfg),
				grpc.Metrics(mtrcs),
				grpc.ServiceHandler(handle),
				grpc.TraceProvider(traceProvider),
			)
			gr.Add(runner.NewGoMicroGrpcServerRunner(cfg.Service.Name+".grpc", grpcServer))

			// prepare a debug server and add it to the group run.
			debugServer, err := debug.Server(
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)
			if err != nil {
				logger.Error().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(runner.NewGolangHttpServerRunner(cfg.Service.Name+".debug", debugServer))

			logger.Warn().Msgf("starting service %s", cfg.Service.Name)
			grResults := gr.Run(ctx)

			if err := runner.ProcessResults(grResults); err != nil {
				logger.Error().Err(err).Msgf("service %s stopped with error", cfg.Service.Name)
				return err
			}
			logger.Warn().Msgf("service %s stopped without error", cfg.Service.Name)
			return nil
		},
	}
}
