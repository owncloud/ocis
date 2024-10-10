package handlers

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

// NewNatsCheck checks the reachability of a nats server.
func NewNatsCheck(natsCluster string, options ...nats.Option) func(context.Context) error {
	return func(ctx context.Context) error {
		n, err := nats.Connect(natsCluster, options...)
		if err != nil {
			return fmt.Errorf("could not connect to nats server: %v", err)
		}
		defer n.Close()
		if n.Status() != nats.CONNECTED {
			return fmt.Errorf("nats server not connected")
		}
		return nil
	}
}
