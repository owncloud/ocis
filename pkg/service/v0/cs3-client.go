package svc

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

func getGatewayServiceClient(endpoint string) (gateway.GatewayAPIClient, error) {
	conn, err := newConn(endpoint)
	if err != nil {
		return nil, err
	}

	return gateway.NewGatewayAPIClient(conn), nil
}

// GetClient returns a gateway client to talk to reva
func (g Graph) GetClient() (gateway.GatewayAPIClient, error) {
	return getGatewayServiceClient(g.config.Reva.Address)
}
