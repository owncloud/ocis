// services/llm/pkg/command/server.go
package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/owncloud/reva/v2/pkg/store"
	"github.com/urfave/cli/v2"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/llm/pkg/config"
	"github.com/owncloud/ocis/v2/services/llm/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/llm/pkg/logging"
	"github.com/owncloud/ocis/v2/services/llm/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/llm/pkg/ratelimit"
	"github.com/owncloud/ocis/v2/services/llm/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/llm/pkg/server/http"
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
			tracerProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to initialize tracer")
				return err
			}

			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			m := metrics.New()
			m.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			rateLimitStore := store.Create(
				store.Store(cfg.Store.Store),
				store.TTL(cfg.Store.TTL),
				microstore.Nodes(cfg.Store.Nodes...),
				microstore.Database(cfg.Store.Database),
				microstore.Table(cfg.Store.Table),
				store.Authentication(cfg.Store.AuthUsername, cfg.Store.AuthPassword),
				store.TLS(cfg.Store.EnableTLS, cfg.Store.TLSInsecure, cfg.Store.TLSRootCACertificate),
			)
			limiter := ratelimit.New(rateLimitStore, cfg.RateLimit.Window, cfg.RateLimit.MaxRequests)

			gr := runner.NewGroup()
			{
				httpServer, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Limiter(limiter),
					http.TraceProvider(tracerProvider),
				)
				if err != nil {
					logger.Info().Err(err).Str("server", "http").Msg("Failed to initialize server")
					return err
				}
				gr.Add(runner.NewGoMicroHttpServerRunner(cfg.Service.Name+".http", httpServer))
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
