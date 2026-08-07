package revaconfig_test

import (
	"testing"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/revaconfig"
)

// TestIDCacheTTLIsPinnedToZero guards the storage-users ID-cache fix: the
// id<->path index must never expire, so every driver config that builds an
// "idcache" block has to pin cache_ttl to 0 regardless of operator config.
// Expiring it makes the provider lose track of existing nodes (files vanish on
// POSIX, re-resolve thrash on the decomposed drivers).
func TestIDCacheTTLIsPinnedToZero(t *testing.T) {
	cfg := defaults.DefaultConfig()
	// The driver builders dereference cfg.Commons.GRPCClientTLS; DefaultConfig
	// does not populate Commons, so set the minimum needed to build the configs.
	cfg.Commons = &shared.Commons{GRPCClientTLS: &shared.GRPCClientTLS{}}

	drivers := map[string]map[string]interface{}{
		"Posix":        revaconfig.Posix(cfg, false),
		"Ocis":         revaconfig.Ocis(cfg),
		"OcisNoEvents": revaconfig.OcisNoEvents(cfg),
		"S3NG":         revaconfig.S3NG(cfg),
		"S3NGNoEvents": revaconfig.S3NGNoEvents(cfg),
	}

	for name, driver := range drivers {
		idcache, ok := driver["idcache"].(map[string]interface{})
		if !ok {
			t.Fatalf("%s: missing idcache config block", name)
		}
		if ttl := idcache["cache_ttl"]; ttl != time.Duration(0) {
			t.Errorf("%s: idcache cache_ttl = %v, want 0 (the id<->path index must not expire)", name, ttl)
		}
	}
}
