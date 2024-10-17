package checks

import (
	"context"
	"fmt"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

// NewGRPCCheck checks the reachability of a grpc server.
func NewGRPCCheck(address string) func(context.Context) error {
	return func(_ context.Context) error {
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("could not connect to grpc server: %v", err)
		}
		_ = conn.Close()
		return nil
	}
}
