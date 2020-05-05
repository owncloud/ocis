package cs3

import (
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"

	"google.golang.org/grpc"
)

func newConn(endpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
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
