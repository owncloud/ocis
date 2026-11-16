package revaconfig

import (
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/sharing/pkg/config/defaults"
)

// TestEventsNatsSettings checks that each events map in the reva config of the
// sharing service gets all the connection, TLS and authentication settings of
// the events broker.
func TestEventsNatsSettings(t *testing.T) {
	cfg := defaults.FullDefaultConfig()
	// The runtime sets the gRPC TLS settings. The defaults do not.
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

	revaCfg, err := SharingConfigFromStruct(cfg, log.NopLogger())
	if err != nil {
		t.Fatal(err)
	}

	grpc := revaCfg["grpc"].(map[string]interface{})
	userShareProvider := grpc["services"].(map[string]interface{})["usershareprovider"].(map[string]interface{})
	jsoncs3 := userShareProvider["drivers"].(map[string]interface{})["jsoncs3"].(map[string]interface{})

	tests := []struct {
		name   string
		events map[string]interface{}
		want   map[string]interface{}
	}{
		{
			name:   "jsoncs3 share manager",
			events: jsoncs3["events"].(map[string]interface{}),
			want: map[string]interface{}{
				"natsaddress":          "nats.example.com:4222",
				"natsclusterid":        "ocis-cluster",
				"enabletls":            true,
				"tlsinsecure":          true,
				"tlsrootcacertificate": "/etc/ocis/ca.crt",
				"authusername":         "nats-user",
				"authpassword":         "nats-password",
				"credsfile":            "/etc/ocis/nats.creds",
				"creds":                "creds-content",
			},
		},
		{
			name:   "events middleware",
			events: grpc["interceptors"].(map[string]interface{})["eventsmiddleware"].(map[string]interface{}),
			want: map[string]interface{}{
				"address":          "nats.example.com:4222",
				"clusterID":        "ocis-cluster",
				"enable-tls":       true,
				"tls-insecure":     true,
				"tls-root-ca-cert": "/etc/ocis/ca.crt",
				"username":         "nats-user",
				"password":         "nats-password",
				"creds-file":       "/etc/ocis/nats.creds",
				"creds":            "creds-content",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.want {
				if got := tt.events[key]; got != value {
					t.Errorf("%s = %v, want %v", key, got, value)
				}
			}
		})
	}
}
