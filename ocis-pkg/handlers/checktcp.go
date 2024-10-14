package handlers

import (
	"context"
	"net"
	"time"
)

// NewTCPCheck returns a check that connects to a given tcp endpoint.
func NewTCPCheck(address string) func(ctx context.Context) error {
	return func(_ context.Context) error {
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil {
			return err
		}

		err = conn.Close()
		if err != nil {
			return err
		}

		return nil
	}
}
