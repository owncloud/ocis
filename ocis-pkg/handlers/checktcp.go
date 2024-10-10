package handlers

import (
	"context"
	"net"
)

// NewTCPCheck returns a check that connects to a given tcp endpoint.
func NewTCPCheck(address string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		conn, err := net.Dial("tcp", address)
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
