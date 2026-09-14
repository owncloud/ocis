package revaconfig

import (
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/config/defaults"
)

func TestStatCacheTLSPlumbing(t *testing.T) {
	t.Setenv("OCIS_CACHE_ENABLE_TLS", "true")
	t.Setenv("OCIS_CACHE_TLS_INSECURE", "true")
	t.Setenv("OCIS_CACHE_TLS_ROOT_CA_CERTIFICATE", "/etc/ocis/ca.crt")

	cfg := defaults.FullDefaultConfig()
	if err := envdecode.Decode(cfg); err != nil {
		t.Fatal(err)
	}

	revaCfg, err := FrontendConfigFromStruct(cfg, log.NopLogger())
	if err != nil {
		t.Fatal(err)
	}

	statCache := revaCfg["http"].(map[string]interface{})["services"].(map[string]interface{})["ocs"].(map[string]interface{})["stat_cache_config"].(map[string]interface{})
	t.Logf("cache_enable_tls=%v cache_tls_insecure=%v cache_tls_root_ca_certificate=%v",
		statCache["cache_enable_tls"], statCache["cache_tls_insecure"], statCache["cache_tls_root_ca_certificate"])

	if statCache["cache_enable_tls"] != true {
		t.Errorf("cache_enable_tls = %v, want true", statCache["cache_enable_tls"])
	}
	if statCache["cache_tls_insecure"] != true {
		t.Errorf("cache_tls_insecure = %v, want true", statCache["cache_tls_insecure"])
	}
	if statCache["cache_tls_root_ca_certificate"] != "/etc/ocis/ca.crt" {
		t.Errorf("cache_tls_root_ca_certificate = %v, want /etc/ocis/ca.crt", statCache["cache_tls_root_ca_certificate"])
	}
}
