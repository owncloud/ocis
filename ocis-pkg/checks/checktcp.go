package checks

import (
	"context"
	"net"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
)

// NewTCPCheck returns a check that connects to a given tcp endpoint.
func NewTCPCheck(address string) func(context.Context) error {
	return func(_ context.Context) error {
		address, err := handlers.FailSaveAddress(address)
		if err != nil {
			return err
		}

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
