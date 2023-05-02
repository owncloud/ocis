package command

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/oklog/run"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/logging"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/service"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
			return err
		},
		Action: func(c *cli.Context) error {
			var (
				gr     = run.Group{}
				logger = logging.Configure(cfg.Service.Name, cfg.Log)

				evtsCfg = cfg.Postprocessing.Events
				tlsConf *tls.Config

				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
			)
			defer cancel()

			{
				if evtsCfg.EnableTLS {
					var rootCAPool *x509.CertPool
					if evtsCfg.TLSRootCACertificate != "" {
						rootCrtFile, err := os.Open(evtsCfg.TLSRootCACertificate)
						if err != nil {
							return err
						}

						rootCAPool, err = ociscrypto.NewCertPoolFromPEM(rootCrtFile)
						if err != nil {
							return err
						}
						evtsCfg.TLSInsecure = false
					}

					tlsConf = &tls.Config{
						RootCAs: rootCAPool,
					}
				}

				bus, err := stream.Nats(
					natsjs.TLSConfig(tlsConf),
					natsjs.Address(evtsCfg.Endpoint),
					natsjs.ClusterID(evtsCfg.Cluster),
				)
				if err != nil {
					return err
				}

				svc, err := service.NewPostprocessingService(bus, logger, cfg.Postprocessing)
				if err != nil {
					return err
				}
				gr.Add(func() error {
					err := make(chan error)
					select {
					case <-ctx.Done():
						return nil

					case err <- svc.Run():
						return <-err
					}
				}, func(err error) {
					logger.Error().
						Err(err).
						Msg("Shutting down server")
					cancel()
				})
			}

			{
				server := debug.NewService(
					debug.Logger(logger),
					debug.Name(cfg.Service.Name),
					debug.Version(version.GetString()),
					debug.Address(cfg.Debug.Addr),
					debug.Token(cfg.Debug.Token),
					debug.Pprof(cfg.Debug.Pprof),
					debug.Zpages(cfg.Debug.Zpages),
					debug.Health(handlers.Health),
					debug.Ready(handlers.Ready),
				)

				gr.Add(server.ListenAndServe, func(_ error) {
					_ = server.Shutdown(ctx)
					cancel()
				})
			}
			return gr.Run()
		},
	}
}
