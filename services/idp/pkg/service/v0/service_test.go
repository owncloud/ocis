package svc

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
	"go.opentelemetry.io/otel/trace"
)

func TestService(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "idp-test-*")
	defer os.RemoveAll(tmpDir)

	// Setup minimal IDP service config
	encryptionKey := filepath.Join(tmpDir, "encryption.key")
	os.WriteFile(encryptionKey, []byte("test-encryption-secret-32-bytes!"), 0600)

	cfg := &config.Config{
		IDP: config.Settings{
			Iss:                      "https://localhost:9200",
			IdentityManager:          "guest",
			Insecure:                 true,
			AllowClientGuests:        true,
			EncryptionSecretFile:     encryptionKey,
			SigningMethod:            "PS256",
			IdentifierRegistrationConf: filepath.Join(tmpDir, "clients.yaml"),
		},
		Clients: []config.Client{{
			ID: "test", Trusted: true,
			RedirectURIs: []string{"https://localhost:9200/callback"},
		}},
		Commons: &shared.Commons{OcisURL: "https://localhost:9200"},
	}

	// Create actual IDP service
	svc := NewService(
		Config(cfg),
		Logger(log.NewLogger(log.Level("error"))),
		TraceProvider(trace.NewNoopTracerProvider()),
	)

	srv := httptest.NewServer(svc)
	defer srv.Close()

	tests := []struct {
		name     string
		path     string
		wantCode int
		contains string
	}{
		{"identifier_route", "/signin/v1/identifier", 200, ""},
		{"identifier_index", "/signin/v1/identifier/index.html", 200, ""},
		{"welcome_route", "/signin/v1/welcome", 200, ""},
		{"goodbye_route", "/signin/v1/goodbye", 200, ""},
		{"oidc_discovery", "/.well-known/openid-configuration", 200, "issuer"},
		{"jwks_endpoint", "/konnect/v1/jwks.json", 200, "keys"},
		{"handles_204", "/signin/v1/identifier", 200, "res.status === 204"},
		{"logon_mode", "/signin/v1/identifier", 200, `params: [username, password, "1"]`},
		{"query_params", "/signin/v1/identifier", 200, "window.location.search"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(srv.URL + tt.path)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantCode)
			}

			if tt.contains != "" {
				body, _ := io.ReadAll(resp.Body)
				if !strings.Contains(string(body), tt.contains) {
					t.Errorf("missing %q", tt.contains)
				}
			}
		})
	}
}
