package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	mgrpcc "github.com/go-micro/plugins/v4/client/grpc"
	mtracer "github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/reva/v2/pkg/autoprop"
	"go-micro.dev/v4/client"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// ClientOptions represent options (e.g. tls settings) for the grpc clients
type ClientOptions struct {
	tlsMode string
	caCert  string
	tp      trace.TracerProvider
	kp      keepalive.ClientParameters
}

// Option is used to pass client options
type ClientOption func(opts *ClientOptions)

// WithTLSMode allows setting the TLSMode option for grpc clients
func WithTLSMode(v string) ClientOption {
	return func(o *ClientOptions) {
		o.tlsMode = v
	}
}

// WithTLSCACert allows setting the CA Certificate for grpc clients
func WithTLSCACert(v string) ClientOption {
	return func(o *ClientOptions) {
		o.caCert = v
	}
}

// WithTraceProvider allows to set the trace Provider for grpc clients
func WithTraceProvider(tp trace.TracerProvider) ClientOption {
	return func(o *ClientOptions) {
		if tp != nil {
			o.tp = tp
		} else {
			o.tp = noop.NewTracerProvider()
		}
	}
}

// WithKeepaliveParams allows setting the grpc keepalive params for grpc clients, derived from
// co.GRPCClientKeepaliveTime/-Timeout. Without a configured keepalive time this is a no-op and the
// clients keep grpc's own behavior, see GetClientKeepaliveParams.
func WithKeepaliveParams(co *shared.GRPCClientOptions) ClientOption {
	return func(o *ClientOptions) {
		if co == nil {
			return
		}
		o.kp = GetClientKeepaliveParams(co.GRPCClientKeepaliveTime, co.GRPCClientKeepaliveTimeout)
	}
}

func GetClientOptions(t *shared.GRPCClientTLS) []ClientOption {
	opts := []ClientOption{
		WithTLSMode(t.Mode),
		WithTLSCACert(t.CACert),
	}
	return opts
}

func NewClient(opts ...ClientOption) (client.Client, error) {
	var options ClientOptions
	for _, opt := range opts {
		opt(&options)
	}

	reg := registry.GetRegistry()
	var tlsConfig *tls.Config
	cOpts := []client.Option{
		client.Registry(reg),
		client.Wrap(mtracer.NewClientWrapper(
			mtracer.WithTraceProvider(options.tp),
		)),
		client.Wrap(autoprop.NewGoMicroClientWrapper()),
	}

	if options.kp != (keepalive.ClientParameters{}) {
		// same mechanism the reva pool uses in grpc.NewConn (pkg/rgrpc/todo/pool/connection.go),
		// so that a black-holed peer is detected here as well instead of hanging for the caller's full timeout
		cOpts = append(cOpts, withDefaultCallOptions(mgrpcc.DialOptions(grpc.WithKeepaliveParams(options.kp))))
	}

	switch options.tlsMode {
	case "insecure":
		if os.Getenv("OCIS_INSECURE") != "true" {
			return nil, errors.New("insecure TLS mode is only allowed in development environments with OCIS_INSECURE=true")
		}
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		cOpts = append(cOpts, mgrpcc.AuthTLS(tlsConfig))
	case "on":
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		// Note: If caCert is empty we use the system's default set of trusted CAs
		if options.caCert != "" {
			certs := x509.NewCertPool()
			pemData, err := os.ReadFile(options.caCert)
			if err != nil {
				return nil, err
			}
			if !certs.AppendCertsFromPEM(pemData) {
				return nil, errors.New("could not initialize client, adding CA cert failed")
			}
			tlsConfig.RootCAs = certs
		}
		cOpts = append(cOpts, mgrpcc.AuthTLS(tlsConfig))
		// case "off":
		// default:
	}

	return mgrpcc.NewClient(cOpts...), nil
}

// withDefaultCallOptions applies client.CallOptions (e.g. mgrpcc.DialOptions) to the client's
// default CallOptions. The go-micro grpc plugin only reads dial options back out of the default
// CallOptions.Context set at construction time (not per-call CallOptions), so this is the only
// way to make them apply to every dial the client makes.
func withDefaultCallOptions(callOpts ...client.CallOption) client.Option {
	return func(o *client.Options) {
		for _, co := range callOpts {
			co(&o.CallOptions)
		}
	}
}
