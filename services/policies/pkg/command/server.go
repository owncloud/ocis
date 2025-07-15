package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	svcProtogen "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine/opa"
	"github.com/owncloud/ocis/v2/services/policies/pkg/server/debug"
	svcEvent "github.com/owncloud/ocis/v2/services/policies/pkg/service/event"
	svcGRPC "github.com/owncloud/ocis/v2/services/policies/pkg/service/grpc"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", "authz"),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			logger := log.NewLogger(
				log.Name(cfg.Service.Name),
				log.Level(cfg.Log.Level),
				log.Pretty(cfg.Log.Pretty),
				log.Color(cfg.Log.Color),
				log.File(cfg.Log.File),
			).SubloggerWithRequestID(ctx)

			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			e, err := opa.NewOPA(cfg.Engine.Timeout, logger, cfg.Engine)
			if err != nil {
				return err
			}

			gr := runner.NewGroup()
			{
				grpcClient, err := grpc.NewClient(
					append(
						grpc.GetClientOptions(cfg.GRPCClientTLS),
						grpc.WithTraceProvider(traceProvider),
					)...,
				)
				if err != nil {
					return err
				}

				svc, err := grpc.NewServiceWithClient(
					grpcClient,
					grpc.Logger(logger),
					grpc.TLSEnabled(cfg.GRPC.TLS.Enabled),
					grpc.TLSCert(
						cfg.GRPC.TLS.Cert,
						cfg.GRPC.TLS.Key,
					),
					grpc.Name(cfg.Service.Name),
					grpc.Context(ctx),
					grpc.Address(cfg.GRPC.Addr),
					grpc.Namespace(cfg.GRPC.Namespace),
					grpc.Version(version.GetString()),
					grpc.TraceProvider(traceProvider),
				)
				if err != nil {
					return err
				}

				grpcSvc, err := svcGRPC.New(e)
				if err != nil {
					return err
				}

				if err := svcProtogen.RegisterPoliciesProviderHandler(
					svc.Server(),
					grpcSvc,
				); err != nil {
					return err
				}

				gr.Add(runner.NewGoMicroGrpcServerRunner(cfg.Service.Name+".grpc", svc))
			}

			{

				connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
				bus, err := stream.NatsFromConfig(connName, false, stream.NatsConfig(cfg.Events))
				if err != nil {
					return err
				}

				eventSvc, err := svcEvent.New(ctx, bus, logger, traceProvider, e, cfg.Postprocessing.Query)
				if err != nil {
					return err
				}

				gr.Add(runner.New(cfg.Service.Name+".svc", func() error {
					return eventSvc.Run()
				}, func() {
					eventSvc.Close()
				}))
			}

			{
				debugServer, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)
				if err != nil {
					logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(runner.NewGolangHttpServerRunner(cfg.Service.Name+".debug", debugServer))
			}

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
