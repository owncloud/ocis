package command

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/owncloud/ocis-glauth/pkg/crypto"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	glauthcfg "github.com/glauth/glauth/pkg/config"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	accounts "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-glauth/pkg/config"
	"github.com/owncloud/ocis-glauth/pkg/flagset"
	"github.com/owncloud/ocis-glauth/pkg/server/debug"
	"github.com/owncloud/ocis-glauth/pkg/server/glauth"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: flagset.ServerWithConfig(cfg),
		Before: func(c *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			return ParseConfig(c, cfg)
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					exporter, err := ocagent.NewExporter(
						ocagent.WithReconnectionPeriod(5*time.Second),
						ocagent.WithAddress(cfg.Tracing.Endpoint),
						ocagent.WithServiceName(cfg.Tracing.Service),
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create agent tracing")

						return err
					}

					trace.RegisterExporter(exporter)
					view.RegisterExporter(exporter)

				case "jaeger":
					exporter, err := jaeger.NewExporter(
						jaeger.Options{
							AgentEndpoint:     cfg.Tracing.Endpoint,
							CollectorEndpoint: cfg.Tracing.Collector,
							ServiceName:       cfg.Tracing.Service,
						},
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create jaeger tracing")

						return err
					}

					trace.RegisterExporter(exporter)

				case "zipkin":
					endpoint, err := openzipkin.NewEndpoint(
						cfg.Tracing.Service,
						cfg.Tracing.Endpoint,
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create zipkin tracing")

						return err
					}

					exporter := zipkin.NewExporter(
						zipkinhttp.NewReporter(
							cfg.Tracing.Collector,
						),
						endpoint,
					)

					trace.RegisterExporter(exporter)

				default:
					logger.Warn().
						Str("type", t).
						Msg("Unknown tracing backend")
				}

				trace.ApplyConfig(
					trace.Config{
						DefaultSampler: trace.AlwaysSample(),
					},
				)
			} else {
				logger.Debug().
					Msg("Tracing is not enabled")
			}

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
				//metrics     = metrics.New()
			)

			defer cancel()

			{
				cfg := glauthcfg.Config{
					LDAP: glauthcfg.LDAP{
						Enabled: cfg.Ldap.Enabled,
						Listen:  cfg.Ldap.Address,
					},
					LDAPS: glauthcfg.LDAPS{
						Enabled: cfg.Ldaps.Enabled,
						Listen:  cfg.Ldaps.Address,
						Cert:    cfg.Ldaps.Cert,
						Key:     cfg.Ldaps.Key,
					},
					Backend: glauthcfg.Backend{
						BaseDN:      cfg.Backend.BaseDN,
						Insecure:    cfg.Backend.Insecure,
						NameFormat:  cfg.Backend.NameFormat,
						GroupFormat: cfg.Backend.GroupFormat,
						SSHKeyAttr:  cfg.Backend.SSHKeyAttr,
					},
				}

				if cfg.LDAPS.Enabled {
					// GenCert has side effects as it writes 2 files to the binary running location
					if err := crypto.GenCert("ldap.crt", "ldap.key", logger); err != nil {
						logger.Fatal().Err(err).Msgf("Could not generate test-certificate")
					}
				}

				as, gs, err := getAccountsServices()
				if err != nil {
					return err
				}

				server, err := glauth.Server(
					glauth.AccountsService(as),
					glauth.GroupsService(gs),
					glauth.Logger(logger),
					glauth.Config(&cfg),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("transport", "ldap").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					err := make(chan error)
					select {
					case <-ctx.Done():
						return nil
					case err <- server.ListenAndServe():
						return <-err
					}

				}, func(_ error) {
					logger.Info().
						Str("transport", "ldap").
						Msg("Shutting down server")

					server.Shutdown()
					cancel()
				})

				gr.Add(func() error {
					err := make(chan error)
					select {
					case <-ctx.Done():
						return nil
					case err <- server.ListenAndServeTLS():
						return <-err
					}

				}, func(_ error) {
					logger.Info().
						Str("transport", "ldaps").
						Msg("Shutting down server")

					server.Shutdown()
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
					logger.Info().
						Err(err).
						Str("transport", "debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)

					defer timeout()
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						logger.Info().
							Err(err).
							Str("transport", "debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("transport", "debug").
							Msg("Shutting down server")
					}
				})
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}

// getAccountsServices returns an ocis-accounts service
func getAccountsServices() (accounts.AccountsService, accounts.GroupsService, error) {
	service := micro.NewService()

	// parse command line flags
	service.Init()

	err := service.Client().Init(
		client.ContentType("application/json"),
	)
	if err != nil {
		return nil, nil, err
	}
	return accounts.NewAccountsService("com.owncloud.api.accounts", service.Client()),
		accounts.NewGroupsService("com.owncloud.api.accounts", service.Client()),
		nil
}
