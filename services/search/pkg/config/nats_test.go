package config

import (
	"reflect"
	"testing"
)

// TestNatsConfigSetsAllFields fails when a new field of stream.NatsConfig is not mapped.
func TestNatsConfigSetsAllFields(t *testing.T) {
	cfg := Events{
		Endpoint:             "nats:9233",
		Cluster:              "ocis-cluster",
		EnableTLS:            true,
		TLSInsecure:          true,
		TLSRootCACertificate: "/ca.pem",
		AuthUsername:         "user",
		AuthPassword:         "pass",
		CredsFile:            "/user.creds",
		Creds:                "creds",
	}
	got := reflect.ValueOf(cfg.NatsConfig())
	for i := 0; i < got.NumField(); i++ {
		name := got.Type().Field(i).Name
		if got.Field(i).IsZero() {
			t.Errorf("field %s is not set", name)
		}
	}
}
