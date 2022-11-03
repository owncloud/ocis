package grpc

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	mgrpcs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	"github.com/go-micro/plugins/v4/wrapper/trace/opencensus"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
	mtls "go-micro.dev/v4/util/tls"
)

// Service simply wraps the go-micro grpc service.
type Service struct {
	micro.Service
}

// NewService initializes a new grpc service.
func NewService(opts ...Option) (Service, error) {
	var mServer server.Server
	sopts := newOptions(opts...)
	tlsConfig := &tls.Config{}
	if sopts.TLSEnabled {
		var cert tls.Certificate
		var err error
		if sopts.TLSCert != "" {
			cert, err = tls.LoadX509KeyPair(sopts.TLSCert, sopts.TLSKey)
			if err != nil {
				sopts.Logger.Error().Err(err).Str("cert", sopts.TLSCert).Str("key", sopts.TLSKey).Msg("error loading server certifcate and key")
				return Service{}, fmt.Errorf("grpc service error loading server certificate and key: %w", err)
			}
		} else {
			// Generate a self-signed server certificate on the fly. This requires the clients
			// to connect with InsecureSkipVerify.
			subj := []string{sopts.Address}
			if host, _, err := net.SplitHostPort(sopts.Address); err == nil && host != "" {
				subj = []string{host}
			}

			sopts.Logger.Warn().Str("address", sopts.Address).
				Msg("GRPC: No server certificate configured. Generating a temporary self-signed certificate")

			cert, err = mtls.Certificate(subj...)
			if err != nil {
				return Service{}, fmt.Errorf("grpc service error creating temporary self-signed certificate: %w", err)
			}
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		mServer = mgrpcs.NewServer(mgrpcs.AuthTLS(tlsConfig))
	} else {
		mServer = mgrpcs.NewServer()
	}

	mopts := []micro.Option{
		// first add a server because it will reset any options
		micro.Server(mServer),
		// also add a client that can be used after initializing the service
		micro.Client(DefaultClient()),
		micro.Address(sopts.Address),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Registry(registry.GetRegistry()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
	}

	return Service{micro.NewService(mopts...)}, nil
}
