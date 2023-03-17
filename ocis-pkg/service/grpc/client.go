package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"sync"

	mgrpcc "github.com/go-micro/plugins/v4/client/grpc"
	mbreaker "github.com/go-micro/plugins/v4/wrapper/breaker/gobreaker"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"go-micro.dev/v4/client"
)

var (
	defaultClient client.Client
	once          sync.Once
)

// ClientOptions represent options (e.g. tls settings) for the grpc clients
type ClientOptions struct {
	tlsMode string
	caCert  string
}

// Option is used to pass client options
type ClientOption func(opts *ClientOptions)

// WithTLSMode allows to set the TLSMode option for grpc clients
func WithTLSMode(v string) ClientOption {
	return func(o *ClientOptions) {
		o.tlsMode = v
	}
}

// WithTLSCACert allows to set the CA Certificate for grpc clients
func WithTLSCACert(v string) ClientOption {
	return func(o *ClientOptions) {
		o.caCert = v
	}
}

// Configure configures the default oOCIS grpc client (e.g. TLS settings)
func Configure(opts ...ClientOption) error {
	var options ClientOptions
	for _, opt := range opts {
		opt(&options)
	}

	var outerr error
	once.Do(func() {
		reg := registry.GetRegistry()
		var tlsConfig *tls.Config
		cOpts := []client.Option{
			client.Registry(reg),
			client.Wrap(mbreaker.NewClientWrapper()),
		}
		switch options.tlsMode {
		case "insecure":
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
			cOpts = append(cOpts, mgrpcc.AuthTLS(tlsConfig))
		case "on":
			// TODO use function to add cert to cert pool?
			tlsConfig = &tls.Config{}
			// Note: If caCert is empty we use the system's default set of trusted CAs
			if options.caCert != "" {
				certs := x509.NewCertPool()
				pemData, err := os.ReadFile(options.caCert)
				if err != nil {
					outerr = err
					return
				}
				if !certs.AppendCertsFromPEM(pemData) {
					outerr = errors.New("Error initializing LDAP Backend. Adding CA cert failed")
					return
				}
				tlsConfig.RootCAs = certs
			}
			cOpts = append(cOpts, mgrpcc.AuthTLS(tlsConfig))
		}

		defaultClient = mgrpcc.NewClient(cOpts...)
	})
	return outerr
}

// DefaultClient returns a custom oCIS grpc configured client.
func DefaultClient() client.Client {
	return defaultClient
}

func GetClientOptions(t *shared.GRPCClientTLS) []ClientOption {
	opts := []ClientOption{
		WithTLSMode(t.Mode),
		WithTLSCACert(t.CACert),
	}
	return opts
}
