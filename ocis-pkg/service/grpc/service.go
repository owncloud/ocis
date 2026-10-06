package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	mgrpcs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	mtracer "github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"github.com/rs/zerolog"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/reva/v2/pkg/autoprop"
)

// _minKeepalivePingInterval is the shortest keepalive ping interval we accept
// from a client. Clients cannot ping more often than every ten seconds, so this
// leaves headroom below that; grpc's own default would GOAWAY anything more
// frequent than every five minutes.
const _minKeepalivePingInterval = 5 * time.Second

// Service simply wraps the go-micro grpc service.
type Service struct {
	micro.Service
}

// NewServiceWithClient initializes a new grpc service with explicit client.
func NewServiceWithClient(client client.Client, opts ...Option) (Service, error) {
	var mServer server.Server
	sopts := newOptions(opts...)
	// accept the keepalive pings clients send while an rpc is in flight, so that
	// they can tell a peer that stopped answering from one that is merely slow
	keepaliveParams := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             _minKeepalivePingInterval,
		PermitWithoutStream: true,
	})
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
			cert, err = ociscrypto.GenTempCertForAddr(sopts.Address)
			if err != nil {
				return Service{}, fmt.Errorf("grpc service error creating temporary self-signed certificate: %w", err)
			}
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		mServer = mgrpcs.NewServer(mgrpcs.Options(keepaliveParams), mgrpcs.AuthTLS(tlsConfig))
	} else {
		mServer = mgrpcs.NewServer(mgrpcs.Options(keepaliveParams))
	}

	handlerWrappers := []server.HandlerWrapper{
		mtracer.NewHandlerWrapper(
			mtracer.WithTraceProvider(sopts.TraceProvider),
		),
		autoprop.NewGoMicroServerHandlerWrapper(),
	}
	if sopts.Logger.GetLevel() == zerolog.DebugLevel {
		handlerWrappers = append(handlerWrappers, LogHandler(&sopts.Logger))
	}
	handlerWrappers = append(handlerWrappers, sopts.HandlerWrappers...)

	mopts := []micro.Option{
		// first add a server because it will reset any options
		micro.Server(mServer),
		// also add a client that can be used after initializing the service
		micro.Client(client),
		micro.Address(sopts.Address),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
		micro.Context(sopts.Context),
		micro.Registry(registry.GetRegistry()),
		micro.RegisterTTL(registry.GetRegisterTTL()),
		micro.RegisterInterval(registry.GetRegisterInterval()),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(
			mtracer.NewClientWrapper(
				mtracer.WithTraceProvider(sopts.TraceProvider),
			),
			autoprop.NewGoMicroClientWrapper(),
		),
		micro.WrapHandler(handlerWrappers...),
		micro.WrapSubscriber(
			mtracer.NewSubscriberWrapper(mtracer.WithTraceProvider(sopts.TraceProvider)),
			autoprop.NewGoMicroServerSubscriberWrapper(),
		),
	}

	return Service{micro.NewService(mopts...)}, nil
}

// If used with tracing, please ensure this is registered (by micro.WrapHandler()) after
// micro-plugin's opentracing wrapper: `opentracing.NewHandlerWrapper()`
func LogHandler(l *log.Logger) func(fn server.HandlerFunc) server.HandlerFunc {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			now := time.Now()
			spanContext := trace.SpanContextFromContext(ctx)
			defer func() {
				l.Debug().
					Str("traceid", spanContext.TraceID().String()).
					Str("method", req.Method()).
					Str("endpoint", req.Endpoint()).
					Str("content-type", req.ContentType()).
					Str("service", req.Service()).
					Interface("headers", req.Header()).
					Dur("duration", time.Since(now)).
					Msg("grpc call")
			}()
			return fn(ctx, req, rsp)
		}
	}
}
