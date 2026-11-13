package config

import "github.com/owncloud/reva/v2/pkg/events/stream"

// ProviderCacheNatsConfig returns the reva NATS config for the provider cache.
// The cache can have more than one node, so the config has no endpoint.
func (c Cache) ProviderCacheNatsConfig() stream.NatsConfig {
	return stream.NatsConfig{
		EnableTLS:            c.ProviderCacheEnableTLS,
		TLSInsecure:          c.ProviderCacheTLSInsecure,
		TLSRootCACertificate: c.ProviderCacheTLSRootCACertificate,
		AuthUsername:         c.ProviderCacheAuthUsername,
		AuthPassword:         c.ProviderCacheAuthPassword,
		CredsFile:            c.ProviderCacheCredsFile,
		Creds:                c.ProviderCacheCreds,
	}
}
