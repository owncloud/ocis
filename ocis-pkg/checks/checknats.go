package checks

import (
	"context"
	"fmt"

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
func NewNatsCheckWithConfig(address string, cfg stream.NatsConfig) func(context.Context) error {
	return func(ctx context.Context) error {
		options, err := NatsOptions(cfg)
		if err != nil {
			return fmt.Errorf("could not configure the nats connection: %v", err)
		}
		return NewNatsCheck(address, options...)(ctx)
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
