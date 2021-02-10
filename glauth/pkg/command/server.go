package command

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/service/grpc"

	"github.com/owncloud/ocis/glauth/pkg/metrics"

	"github.com/owncloud/ocis/glauth/pkg/crypto"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	glauthcfg "github.com/glauth/glauth/pkg/config"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/glauth/pkg/config"
	"github.com/owncloud/ocis/glauth/pkg/flagset"
	"github.com/owncloud/ocis/glauth/pkg/server/debug"
	"github.com/owncloud/ocis/glauth/pkg/server/glauth"
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

			cfg.Backend.Servers = c.StringSlice("backend-server")
			cfg.Fallback.Servers = c.StringSlice("fallback-server")

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
							Process: jaeger.Process{
								ServiceName: cfg.Tracing.Service,
							},
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
				metrics     = metrics.New()
			)

			defer cancel()

			metrics.BuildInfo.WithLabelValues(cfg.Version).Set(1)

			{
				lcfg := glauthcfg.LDAP{
					Enabled: cfg.Ldap.Enabled,
					Listen:  cfg.Ldap.Address,
				}
				lscfg := glauthcfg.LDAPS{
					Enabled: cfg.Ldaps.Enabled,
					Listen:  cfg.Ldaps.Address,
					Cert:    cfg.Ldaps.Cert,
					Key:     cfg.Ldaps.Key,
				}
				bcfg := glauthcfg.Config{
					LDAP:  lcfg,  // TODO remove LDAP from the backend config upstream
					LDAPS: lscfg, // TODO remove LDAP from the backend config upstream
					Backend: glauthcfg.Backend{
						Datastore:   cfg.Backend.Datastore,
						BaseDN:      cfg.Backend.BaseDN,
						Insecure:    cfg.Backend.Insecure,
						NameFormat:  cfg.Backend.NameFormat,
						GroupFormat: cfg.Backend.GroupFormat,
						Servers:     cfg.Backend.Servers,
						SSHKeyAttr:  cfg.Backend.SSHKeyAttr,
						UseGraphAPI: cfg.Backend.UseGraphAPI,
					},
				}
				fcfg := glauthcfg.Config{
					LDAP:  lcfg,  // TODO remove LDAP from the backend config upstream
					LDAPS: lscfg, // TODO remove LDAP from the backend config upstream
					Backend: glauthcfg.Backend{
						Datastore:   cfg.Fallback.Datastore,
						BaseDN:      cfg.Fallback.BaseDN,
						Insecure:    cfg.Fallback.Insecure,
						NameFormat:  cfg.Fallback.NameFormat,
						GroupFormat: cfg.Fallback.GroupFormat,
						Servers:     cfg.Fallback.Servers,
						SSHKeyAttr:  cfg.Fallback.SSHKeyAttr,
						UseGraphAPI: cfg.Fallback.UseGraphAPI,
					},
				}

				if lscfg.Enabled {
					if err := crypto.GenCert(cfg.Ldaps.Cert, cfg.Ldaps.Key, logger); err != nil {
						logger.Fatal().Err(err).Msgf("Could not generate test-certificate")
					}
				}

				as, gs := getAccountsServices()
				server, err := glauth.Server(
					glauth.AccountsService(as),
					glauth.GroupsService(gs),
					glauth.Logger(logger),
					glauth.LDAP(&lcfg),
					glauth.LDAPS(&lscfg),
					glauth.Backend(&bcfg),
					glauth.Fallback(&fcfg),
					glauth.RoleBundleUUID(cfg.RoleBundleUUID),
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
func getAccountsServices() (accounts.AccountsService, accounts.GroupsService) {
	return accounts.NewAccountsService("com.owncloud.api.accounts", grpc.DefaultClient),
		accounts.NewGroupsService("com.owncloud.api.accounts", grpc.DefaultClient)
}
