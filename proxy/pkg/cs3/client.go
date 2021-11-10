package cs3

import (
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	proxytracing "github.com/owncloud/ocis/proxy/pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func newConn(endpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(
				otelgrpc.WithTracerProvider(
					proxytracing.TraceProvider,
				),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// GetGatewayServiceClient returns a new cs3 gateway client
func GetGatewayServiceClient(endpoint string) (gateway.GatewayAPIClient, error) {
	conn, err := newConn(endpoint)
	if err != nil {
		return nil, err
	}

	return gateway.NewGatewayAPIClient(conn), nil
}
