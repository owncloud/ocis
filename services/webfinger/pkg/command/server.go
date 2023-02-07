package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/config"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/logging"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/relations"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/server/http"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/tracing"
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
			err := tracing.Configure(cfg)
			if err != nil {
				return err
			}

			var (
				gr          = run.Group{}
				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
				metrics = metrics.New(metrics.Logger(logger))
			)

			defer cancel()

			metrics.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			{
				svc, err := service.New(
					service.Logger(logger),
					service.Config(cfg),
					service.WithLookupChain(getLookupChain(cfg)),
				)
				if err != nil {
					logger.Error().Err(err).Msg("handler init")
					return err
				}
				svc = service.NewInstrument(svc, metrics)
				svc = service.NewLogging(svc, logger) // this logs service specific data
				svc = service.NewTracing(svc)

				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Service(svc),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("transport", "http").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(err error) {
					logger.Error().
						Err(err).
						Str("transport", "http").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(server.ListenAndServe, func(err error) {
					logger.Error().Err(err)
					_ = server.Shutdown(ctx)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}

func getLookupChain(cfg *config.Config) service.Webfinger {
	lookups := strings.Split(cfg.LookupChain, ",")
	if len(lookups) == 0 {
		return nil
	}
	var webfinger service.Webfinger
	for i := len(lookups) - 1; i >= 0; i-- {
		switch lookups[i] {
		case "openid-discovery":
			webfinger = relations.OpenIDDiscovery(cfg.IDP, webfinger)
		case "owncloud-status":
		case "owncloud-account":
			//url, _ := url.Parse(cfg.OcisURL)
			// TODO error / ignore
			//webfinger = relations.OwnCloudAccount(*url, webfinger)
		case "owncloud-instance":
			webfinger = relations.OwnCloudInstance(cfg.Instances, webfinger)
		default:
			// TODO error / ignore
		}
	}
	return webfinger
}
