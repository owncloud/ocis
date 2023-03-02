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
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	svcProtogen "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
	svcEvent "github.com/owncloud/ocis/v2/services/policies/pkg/service/event"
	svcGRPC "github.com/owncloud/ocis/v2/services/policies/pkg/service/grpc"
	"github.com/urfave/cli/v2"
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
			var (
				gr          = run.Group{}
				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
				logger = log.NewLogger(
					log.Name(cfg.Service.Name),
					log.Level(cfg.Log.Level),
					log.Pretty(cfg.Log.Pretty),
					log.Color(cfg.Log.Color),
					log.File(cfg.Log.File),
				)
			)
			defer cancel()

			e, err := engine.NewOPA(cfg.Engine.Timeout, cfg.Engine)
			if err != nil {
				return err
			}

			{
				err = grpc.Configure(grpc.GetClientOptions(cfg.GRPCClientTLS)...)
				if err != nil {
					return err
				}

				svc, err := grpc.NewService(
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

				gr.Add(svc.Run, func(_ error) {
					cancel()
				})
			}

			{
				var tlsConf *tls.Config

				if cfg.Events.EnableTLS {
					var rootCAPool *x509.CertPool
					if cfg.Events.TLSRootCACertificate != "" {
						rootCrtFile, err := os.Open(cfg.Events.TLSRootCACertificate)
						if err != nil {
							return err
						}

						rootCAPool, err = ociscrypto.NewCertPoolFromPEM(rootCrtFile)
						if err != nil {
							return err
						}
						cfg.Events.TLSInsecure = false
					}

					tlsConf = &tls.Config{
						RootCAs: rootCAPool,
					}
				}

				bus, err := stream.Nats(
					natsjs.TLSConfig(tlsConf),
					natsjs.Address(cfg.Events.Endpoint),
					natsjs.ClusterID(cfg.Events.Cluster),
				)
				if err != nil {
					return err
				}

				eventSvc, err := svcEvent.New(bus, logger, e, cfg.Postprocessing.Query)
				if err != nil {
					return err
				}

				gr.Add(eventSvc.Run, func(_ error) {
					cancel()
				})
			}

			return gr.Run()
		},
	}
}
