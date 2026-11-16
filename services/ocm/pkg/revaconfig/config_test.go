package revaconfig

import (
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/ocm/pkg/config/defaults"
)

// TestEventsNatsSettings checks that each events map in the reva config of the
// ocm service gets all the connection, TLS and authentication settings of the
// events broker.
func TestEventsNatsSettings(t *testing.T) {
	cfg := defaults.FullDefaultConfig()
	// The runtime sets the shared settings and the gRPC TLS settings. The defaults do not.
	cfg.Commons = &shared.Commons{OcisURL: "https://ocis.example.com"}
	cfg.GRPC.TLS = &shared.GRPCServiceTLS{}
	cfg.Events.Endpoint = "nats.example.com:4222"
	cfg.Events.Cluster = "ocis-cluster"
	cfg.Events.EnableTLS = true
	cfg.Events.TLSInsecure = true
	cfg.Events.TLSRootCACertificate = "/etc/ocis/ca.crt"
	cfg.Events.AuthUsername = "nats-user"
	cfg.Events.AuthPassword = "nats-password"
	cfg.Events.CredsFile = "/etc/ocis/nats.creds"
	cfg.Events.Creds = "creds-content"

	revaCfg := OCMConfigFromStruct(cfg, log.NopLogger())

	httpServices := revaCfg["http"].(map[string]interface{})["services"].(map[string]interface{})
	grpcServices := revaCfg["grpc"].(map[string]interface{})["services"].(map[string]interface{})

	want := map[string]interface{}{
		"natsaddress":          "nats.example.com:4222",
		"natsclusterid":        "ocis-cluster",
		"enabletls":            true,
		"tlsinsecure":          true,
		"tlsrootcacertificate": "/etc/ocis/ca.crt",
		"authusername":         "nats-user",
		"authpassword":         "nats-password",
		"credsfile":            "/etc/ocis/nats.creds",
		"creds":                "creds-content",
	}

	tests := []struct {
		name   string
		events map[string]interface{}
	}{
		{
			name:   "sciencemesh",
			events: httpServices["sciencemesh"].(map[string]interface{})["events"].(map[string]interface{}),
		},
		{
			name:   "ocmcore",
			events: grpcServices["ocmcore"].(map[string]interface{})["events"].(map[string]interface{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range want {
				if got := tt.events[key]; got != value {
					t.Errorf("%s = %v, want %v", key, got, value)
				}
			}
		})
	}
}
