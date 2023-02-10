package event

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/go-micro/plugins/v4/events/natsjs"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"go-micro.dev/v4/events"
)

// NewStream prepares the requested nats stream and returns it.
func NewStream(cfg config.Events) (events.Stream, error) {
	var tlsConf *tls.Config

	if cfg.EnableTLS {
		var rootCAPool *x509.CertPool
		if cfg.TLSRootCaCertPath != "" {
			rootCrtFile, err := os.Open(cfg.TLSRootCaCertPath)
			if err != nil {
				return nil, err
			}

			rootCAPool, err = ociscrypto.NewCertPoolFromPEM(rootCrtFile)
			if err != nil {
				return nil, err
			}
			cfg.TLSInsecure = false
		}

		tlsConf = &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    rootCAPool,
		}
	}

	s, err := stream.Nats(
		natsjs.TLSConfig(tlsConf),
		natsjs.Address(cfg.Addr),
		natsjs.ClusterID(cfg.ClusterID),
	)

	if err != nil {
		return nil, err
	}

	return s, nil
}
