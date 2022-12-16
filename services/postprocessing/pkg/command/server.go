package command

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/go-micro/plugins/v4/events/natsjs"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
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
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			evtsCfg := cfg.Postprocessing.Events
			var tlsConf *tls.Config

			if !evtsCfg.TLSInsecure {
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

			bus, err := server.NewNatsStream(
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
			return svc.Run()
		},
	}
}
