package config

import "github.com/owncloud/reva/v2/pkg/events/stream"

// NatsConfig returns the reva NATS config for the events broker.
func (e Events) NatsConfig() stream.NatsConfig {
	return stream.NatsConfig{
		Endpoint:             e.Addr,
		Cluster:              e.ClusterID,
		EnableTLS:            e.EnableTLS,
		TLSInsecure:          e.TLSInsecure,
		TLSRootCACertificate: e.TLSRootCaCertPath,
		AuthUsername:         e.AuthUsername,
		AuthPassword:         e.AuthPassword,
		CredsFile:            e.CredsFile,
		Creds:                e.Creds,
	}
}
