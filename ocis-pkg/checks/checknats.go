package checks

import (
	"context"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/owncloud/reva/v2/pkg/natsauth"
	revastore "github.com/owncloud/reva/v2/pkg/store"
)

// NewNatsCheck checks the reachability of a nats server.
func NewNatsCheck(natsCluster string, options ...nats.Option) func(context.Context) error {
	return func(_ context.Context) error {
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

// NewNatsCheckWithConfig checks the reachability of the nats server at address.
// It authenticates and uses TLS like the service connections configured by cfg.
// The address of cfg is not used.
//
// The check keeps one connection and uses it again for each probe. A new connection
// for each probe reads the .creds file each time. A probe can then read a rotated file
// before the write is complete. Each connection also counts against the connection
// limit of the NATS account. The check connects again only after the connection closes.
// The connection reads the .creds file again on each reconnect.
func NewNatsCheckWithConfig(address string, cfg stream.NatsConfig) func(context.Context) error {
	var (
		mu   sync.Mutex
		conn *nats.Conn
	)
	return func(_ context.Context) error {
		mu.Lock()
		defer mu.Unlock()

		if conn == nil || conn.IsClosed() {
			options, err := NatsOptions(cfg)
			if err != nil {
				return fmt.Errorf("could not configure the nats connection: %v", err)
			}
			c, err := nats.Connect(address, options...)
			if err != nil {
				return fmt.Errorf("could not connect to nats server: %v", err)
			}
			conn = c
		}
		if conn.Status() != nats.CONNECTED {
			return fmt.Errorf("nats server not connected")
		}
		return nil
	}
}

// NatsOptions returns the nats.Option values for the authentication and the TLS settings of cfg.
// JWT credentials take precedence over the username and the password.
func NatsOptions(cfg stream.NatsConfig) ([]nats.Option, error) {
	var options []nats.Option

	credsOpt, err := natsauth.CredentialsOption(cfg.CredsFile, cfg.Creds)
	if err != nil {
		return nil, err
	}
	switch {
	case credsOpt != nil:
		options = append(options, credsOpt)
	case cfg.AuthUsername != "" || cfg.AuthPassword != "":
		options = append(options, nats.UserInfo(cfg.AuthUsername, cfg.AuthPassword))
	}

	if cfg.EnableTLS {
		options = append(options, nats.Secure(revastore.BuildNatsTLSConfig(cfg.TLSInsecure, cfg.TLSRootCACertificate)))
	}
	return options, nil
}
