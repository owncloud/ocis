package command

import (
	"context"
	"github.com/owncloud/ocis-glauth/pkg/crypto"
	"os"
	"os/signal"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	glauthcfg "github.com/glauth/glauth/pkg/config"
	glauth "github.com/glauth/glauth/pkg/server"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/owncloud/ocis-glauth/pkg/config"
	"github.com/owncloud/ocis-glauth/pkg/flagset"
	"github.com/owncloud/ocis-glauth/pkg/mlogr"
	"github.com/owncloud/ocis-glauth/pkg/server/debug"
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
				log := mlogr.New(&logger)
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
						Datastore:   cfg.Backend.Datastore,
						BaseDN:      cfg.Backend.BaseDN,
						Insecure:    cfg.Backend.Insecure,
						NameFormat:  cfg.Backend.NameFormat,
						GroupFormat: cfg.Backend.GroupFormat,
						Servers:     cfg.Backend.Servers,
						SSHKeyAttr:  cfg.Backend.SSHKeyAttr,
						UseGraphAPI: cfg.Backend.UseGraphAPI,
					},
					// TODO read users for the config backend from config file
					Users: []glauthcfg.User{
						glauthcfg.User{
							Name:         "einstein",
							GivenName:    "Albert",
							SN:           "Einstein",
							UnixID:       20000,
							PrimaryGroup: 30000,
							OtherGroups:  []int{30001, 30002, 30007},
							Mail:         "einstein@example.org",
							PassSHA256:   "69bf3575281a970f46e37ecd28b79cfbee6a46e55c10dc91dd36a43410387ab8", // relativity
						},
						glauthcfg.User{
							Name:         "marie",
							GivenName:    "Marie",
							SN:           "Curie",
							UnixID:       20001,
							PrimaryGroup: 30000,
							OtherGroups:  []int{30003, 30004, 30007},
							Mail:         "marie@example.org",
							PassSHA256:   "149a807f82e22b796942efa1010063f4a278cf078ff56ef1d3fc6c156037cef9", // radioactivity
						},
						glauthcfg.User{
							Name:         "feynman",
							GivenName:    "Richard",
							SN:           "Feynman",
							UnixID:       20002,
							PrimaryGroup: 30000,
							OtherGroups:  []int{30005, 30006, 30007},
							Mail:         "feynman@example.org",
							PassSHA256:   "1e2183d3a6017bb01131e27204bb66d3c5fa273acf421c8f9bd4bd633e3d70a8", // superfluidity
						},

						// technical users for ocis
						glauthcfg.User{
							Name:         "konnectd",
							UnixID:       10000,
							PrimaryGroup: 15000,
							Mail:         "idp@example.org",
							PassSHA256:   "e1b6c4460fda166b70f77093f8a2f9b9e0055a5141ed8c6a67cf1105b1af23ca", // konnectd
						},
						glauthcfg.User{
							Name:         "reva",
							UnixID:       10001,
							PrimaryGroup: 15000,
							Mail:         "storage@example.org",
							PassSHA256:   "60a43483d1a41327e689c3ba0451c42661d6a101151e041aa09206305c83e74b", // reva
						},
					},
					Groups: []glauthcfg.Group{
						glauthcfg.Group{
							Name:   "users",
							UnixID: 30000,
						},
						glauthcfg.Group{
							Name:   "sailing-lovers",
							UnixID: 30001,
						},
						glauthcfg.Group{
							Name:   "violin-haters",
							UnixID: 30002,
						},
						glauthcfg.Group{
							Name:   "radium-lovers",
							UnixID: 30003,
						},
						glauthcfg.Group{
							Name:   "polonium-lovers",
							UnixID: 30004,
						},
						glauthcfg.Group{
							Name:   "quantum-lovers",
							UnixID: 30005,
						},
						glauthcfg.Group{
							Name:   "philosophy-haters",
							UnixID: 30006,
						},
						glauthcfg.Group{
							Name:   "physics-lovers",
							UnixID: 30007,
						},
						glauthcfg.Group{
							Name:   "sysusers",
							UnixID: 15000,
						},
					},
				}

				if cfg.LDAPS.Enabled {
					// GenCert has side effects as it writes 2 files to the binary running location
					if err := crypto.GenCert("ldap.crt", "ldap.key", logger); err != nil {
						logger.Fatal().Err(err).Msgf("Could not generate test-certificate")
					}
				}

				server, err := glauth.NewServer(
					glauth.Logger(log),
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
