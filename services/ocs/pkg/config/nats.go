package config

import "github.com/owncloud/reva/v2/pkg/events/stream"

// NatsConfig returns the reva NATS config for the signing keys store.
// The store can have more than one node, so the config has no endpoint.
func (k SigningKeys) NatsConfig() stream.NatsConfig {
	return stream.NatsConfig{
		EnableTLS:            k.EnableTLS,
		TLSInsecure:          k.TLSInsecure,
		TLSRootCACertificate: k.TLSRootCACertificate,
		AuthUsername:         k.AuthUsername,
		AuthPassword:         k.AuthPassword,
		CredsFile:            k.CredsFile,
		Creds:                k.Creds,
	}
}
