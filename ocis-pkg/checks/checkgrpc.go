package checks

import (
	"context"
	"fmt"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewGRPCCheck checks the reachability of a grpc server.
func NewGRPCCheck(address string) func(context.Context) error {
	return func(_ context.Context) error {
		if strings.Contains(address, "0.0.0.0") {
			outboundIp, err := handlers.GetOutBoundIP()
			if err != nil {
				return err
			}
			address = strings.Replace(address, "0.0.0.0", outboundIp, 1)
		}

		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("could not connect to grpc server: %v", err)
		}
		_ = conn.Close()
		return nil
	}
}
