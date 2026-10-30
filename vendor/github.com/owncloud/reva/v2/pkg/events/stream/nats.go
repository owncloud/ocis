package stream

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"os"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/nats-io/nats.go"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/logger"
	"github.com/owncloud/reva/v2/pkg/natsauth"
)

// NatsConfig is the configuration needed for a NATS event stream
type NatsConfig struct {
	Endpoint             string `mapstructure:"address"`          // Endpoint of the nats server
	Cluster              string `mapstructure:"clusterID"`        // CluserID of the nats cluster
	TLSInsecure          bool   `mapstructure:"tls-insecure"`     // Whether to verify TLS certificates
	TLSRootCACertificate string `mapstructure:"tls-root-ca-cert"` // The root CA certificate used to validate the TLS certificate
	EnableTLS            bool   `mapstructure:"enable-tls"`       // Enable TLS
	AuthUsername         string `mapstructure:"username"`         // Username for authentication
	AuthPassword         string `mapstructure:"password"`         // Password for authentication
	CredsFile            string `mapstructure:"creds-file"`       // Path to a NATS .creds file for JWT authentication
	Creds                string `mapstructure:"creds"`            // Content of a NATS .creds file for JWT authentication. Takes precedence over CredsFile
}

// NatsFromConfig returns a nats stream from the given config
func NatsFromConfig(connName string, disableDurability bool, cfg NatsConfig) (events.Stream, error) {
	natsOpts, err := NatsOptions(cfg)
	if err != nil {
		return nil, err
	}
	if (cfg.Creds != "" || cfg.CredsFile != "") && (cfg.AuthUsername != "" || cfg.AuthPassword != "") {
		logger.New().Warn().Str("connection", connName).Msg("nats jwt credentials are configured, ignoring username and password")
	}

	opts := []natsjs.Option{
		natsjs.Address(cfg.Endpoint),
		natsjs.ClusterID(cfg.Cluster),
		natsjs.SynchronousPublish(true),
		natsjs.Name(connName),
		natsjs.NatsOptions(natsOpts...),
	}

	if disableDurability {
		opts = append(opts, natsjs.DisableDurableStreams())
	}

	return Nats(opts...)
}

// NatsOptions returns the nats.Option values for the authentication and the TLS
// settings of cfg. The result does not contain the endpoint and the cluster of cfg.
// JWT credentials take precedence over the username and the password.
// NatsOptions returns an error when it cannot read the credentials or the root CA certificate.
func NatsOptions(cfg NatsConfig) ([]nats.Option, error) {
	var opts []nats.Option

	credsOpt, err := natsauth.CredentialsOption(cfg.CredsFile, cfg.Creds)
	if err != nil {
		return nil, err
	}
	switch {
	case credsOpt != nil:
		opts = append(opts, credsOpt)
	case cfg.AuthUsername != "" || cfg.AuthPassword != "":
		opts = append(opts, nats.UserInfo(cfg.AuthUsername, cfg.AuthPassword))
	}

	if cfg.EnableTLS {
		tlsConf, err := natsTLSConfig(cfg)
		if err != nil {
			return nil, err
		}
		opts = append(opts, nats.Secure(tlsConf))
	}
	return opts, nil
}

// natsTLSConfig returns the TLS configuration of cfg. A root CA certificate turns on
// the verification of the server certificate.
func natsTLSConfig(cfg NatsConfig) (*tls.Config, error) {
	var rootCAPool *x509.CertPool
	if cfg.TLSRootCACertificate != "" {
		rootCrtFile, err := os.Open(cfg.TLSRootCACertificate)
		if err != nil {
			return nil, err
		}
		defer rootCrtFile.Close()

		rootCAPool, err = newCertPoolFromPEM(rootCrtFile)
		if err != nil {
			return nil, err
		}
		cfg.TLSInsecure = false
	}

	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.TLSInsecure, //nolint:gosec
		RootCAs:            rootCAPool,
	}, nil
}

// nats returns a nats streaming client
// retries exponentially to connect to a nats server
func Nats(opts ...natsjs.Option) (events.Stream, error) {
	b := backoff.NewExponentialBackOff()
	var stream events.Stream
	o := func() error {
		n := b.NextBackOff()
		s, err := natsjs.NewStream(opts...)
		if err != nil && n > time.Second {
			logger.New().Error().Err(err).Msgf("can't connect to nats (jetstream) server, retrying in %s", n)
		}
		stream = s
		return err
	}

	err := backoff.Retry(o, b)
	return stream, err
}

// newCertPoolFromPEM reads certificates from io.Reader and returns a x509.CertPool
// containing those certificates.
func newCertPoolFromPEM(crts ...io.Reader) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	var buf bytes.Buffer
	for _, c := range crts {
		if _, err := io.Copy(&buf, c); err != nil {
			return nil, err
		}
		if !certPool.AppendCertsFromPEM(buf.Bytes()) {
			return nil, errors.New("failed to append cert from PEM")
		}
		buf.Reset()
	}

	return certPool, nil
}
