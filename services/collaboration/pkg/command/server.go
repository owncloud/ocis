package command

import (
	"context"
	"fmt"
	"net"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/cs3wopiserver"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/internal/logging"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/server/grpc"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/server/http"
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

			gr := run.Group{}
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
			}()
			defer cancel()

			app, err := cs3wopiserver.Start(cfg) // grpc server needs decoupling
			if err != nil {
				return err
			}
			grpcServer, teardown, err := grpc.Server(
				grpc.App(app),
				grpc.Config(cfg),
				grpc.Logger(logger),
			)
			defer teardown()
			if err != nil {
				logger.Info().
					Err(err).
					Str("transport", "grpc").
					Msg("Failed to initialize server")
				return err
			}

			gr.Add(func() error {
				l, err := net.Listen("tcp", cfg.GRPC.Addr)
				if err != nil {
					return err
				}
				grpcServer.Serve(l)
				return nil
			},
				func(_ error) {
					logger.Error().
						Err(err).
						Str("server", "grpc").
						Msg("shutting down server")
					cancel()
				})

			/*
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)
				if err != nil {
					logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(server.ListenAndServe, func(_ error) {
					_ = server.Shutdown(ctx)
					cancel()
				})
			*/
			server, err := http.Server(
				http.App(app),
				http.Logger(logger),
				http.Config(cfg),
				http.Context(ctx),
				http.TracerProvider(traceProvider),
			)
			gr.Add(server.Run, func(_ error) {
				cancel()
			})

			return gr.Run()
		},
	}
}
