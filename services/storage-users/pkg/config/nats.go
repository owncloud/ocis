package config

import "github.com/owncloud/reva/v2/pkg/events/stream"

// NatsConfig returns the reva NATS config for the events broker.
// Events has the NumConsumers field, so a conversion to stream.NatsConfig is not possible.
func (e Events) NatsConfig() stream.NatsConfig {
	return stream.NatsConfig{
		Endpoint:             e.Endpoint,
		Cluster:              e.Cluster,
		EnableTLS:            e.EnableTLS,
		TLSInsecure:          e.TLSInsecure,
		TLSRootCACertificate: e.TLSRootCACertificate,
		AuthUsername:         e.AuthUsername,
		AuthPassword:         e.AuthPassword,
		CredsFile:            e.CredsFile,
		Creds:                e.Creds,
	}
}
