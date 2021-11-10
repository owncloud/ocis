package cs3

import (
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	proxytracing "github.com/owncloud/ocis/proxy/pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func newConn(endpoint string, insecure bool) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}

	opts = append(opts, grpc.WithUnaryInterceptor(
		otelgrpc.UnaryClientInterceptor(
			otelgrpc.WithTracerProvider(
				proxytracing.TraceProvider,
			),
		),
	))

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(
		endpoint,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// GetGatewayServiceClient returns a new cs3 gateway client
func GetGatewayServiceClient(endpoint string, insecure bool) (gateway.GatewayAPIClient, error) {
	conn, err := newConn(endpoint, insecure)
	if err != nil {
		return nil, err
	}

	return gateway.NewGatewayAPIClient(conn), nil
}
