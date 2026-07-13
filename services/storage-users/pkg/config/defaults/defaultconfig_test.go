package defaults

import (
	"testing"
	"time"
)

// TestDefaultConfigFilemetadataCacheKeepsTTL guards the file-metadata cache TTL,
// which - unlike the ID cache - is a real cache and must keep expiring. The ID
// cache has no TTL knob at all anymore (see the revaconfig test that asserts its
// cache_ttl is pinned to 0); this test ensures that removal did not accidentally
// drop the metadata cache TTL too.
func TestDefaultConfigFilemetadataCacheKeepsTTL(t *testing.T) {
	cfg := DefaultConfig()

	if want := 24 * 60 * time.Second; cfg.FilemetadataCache.TTL != want {
		t.Errorf("FilemetadataCache.TTL = %s, want %s (metadata cache TTL must remain)", cfg.FilemetadataCache.TTL, want)
	}
}
