package config

import (
	"reflect"
	"testing"
)

// TestNatsConfigSetsAllFields fails when a new field of stream.NatsConfig is not mapped.
func TestNatsConfigSetsAllFields(t *testing.T) {
	cfg := Cache{
		ProviderCacheEnableTLS:            true,
		ProviderCacheTLSInsecure:          true,
		ProviderCacheTLSRootCACertificate: "/ca.pem",
		ProviderCacheAuthUsername:         "user",
		ProviderCacheAuthPassword:         "pass",
		ProviderCacheCredsFile:            "/user.creds",
		ProviderCacheCreds:                "creds",
	}
	got := reflect.ValueOf(cfg.ProviderCacheNatsConfig())
	for i := 0; i < got.NumField(); i++ {
		name := got.Type().Field(i).Name
		if name != "Endpoint" && name != "Cluster" && got.Field(i).IsZero() {
			t.Errorf("field %s is not set", name)
		}
	}
}
