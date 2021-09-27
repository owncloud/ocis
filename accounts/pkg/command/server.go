package command

import (
	"context"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/flagset"
	"github.com/owncloud/ocis/accounts/pkg/metrics"
	"github.com/owncloud/ocis/accounts/pkg/server/grpc"
	"github.com/owncloud/ocis/accounts/pkg/server/http"
	svc "github.com/owncloud/ocis/accounts/pkg/service/v0"
	"github.com/owncloud/ocis/accounts/pkg/tracing"
	"github.com/urfave/cli/v2"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start ocis accounts service",
		Description: "uses an LDAP server as the storage backend",
		Flags:       flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			logger := NewLogger(cfg)
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			// When running on single binary mode the before hook from the root command won't get called. We manually
			// call this before hook from ocis command, so the configuration can be loaded.
			if !cfg.Supervised {
				return ParseConfig(ctx, cfg)
			}
			logger.Debug().Str("service", "accounts").Msg("ignoring config file parsing when running supervised")
			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			err := tracing.Configure(cfg)
			if err != nil {
				return err
			}
			gr := run.Group{}
			ctx, cancel := defineContext(cfg)
			mtrcs := metrics.New()

			defer cancel()

			mtrcs.BuildInfo.WithLabelValues(cfg.Server.Version).Set(1)

			handler, err := svc.New(svc.Logger(logger), svc.Config(cfg))
			if err != nil {
				logger.Error().Err(err).Msg("handler init")
				return err
			}

			httpServer := http.Server(
				http.Config(cfg),
				http.Logger(logger),
				http.Name(cfg.Server.Name),
				http.Context(ctx),
				http.Metrics(mtrcs),
				http.Handler(handler),
			)

			gr.Add(httpServer.Run, func(_ error) {
				logger.Info().Str("server", "http").Msg("shutting down server")
				cancel()
			})

			grpcServer := grpc.Server(
				grpc.Config(cfg),
				grpc.Logger(logger),
				grpc.Name(cfg.Server.Name),
				grpc.Context(ctx),
				grpc.Metrics(mtrcs),
				grpc.Handler(handler),
			)

			gr.Add(grpcServer.Run, func(_ error) {
				logger.Info().Str("server", "grpc").Msg("shutting down server")
				cancel()
			})

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// defineContext sets the context for the extension. If there is a context configured it will create a new child from it,
// if not, it will create a root context that can be cancelled.
func defineContext(cfg *config.Config) (context.Context, context.CancelFunc) {
	return func() (context.Context, context.CancelFunc) {
		if cfg.Context == nil {
			return context.WithCancel(context.Background())
		}
		return context.WithCancel(cfg.Context)
	}()
}
