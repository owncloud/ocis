package config

import (
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
)

// TestCacheCredsScopePrecedence pins the precedence against the real Cache
// struct, not only envdecode's synthetic test structs.
func TestCacheCredsScopePrecedence(t *testing.T) {
	t.Setenv("OCIS_NATS_CREDS_FILE", "/global.creds")
	t.Setenv("PROXY_OIDC_USERINFO_CACHE_CREDS", "service-content")

	var c Cache
	if err := envdecode.Decode(&c); err != nil {
		t.Fatal(err)
	}
	if c.CredsFile != "/global.creds" {
		t.Errorf("got file %q, want %q", c.CredsFile, "/global.creds")
	}
	if c.Creds != "service-content" {
		t.Errorf("got creds %q, want %q", c.Creds, "service-content")
	}
}
