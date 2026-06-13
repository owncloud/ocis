package defaults

import (
	"testing"
	"time"
)

// TestDefaultConfigIDCacheHasNoTTL guards against re-introducing a TTL on the
// storage-users ID cache. The ID cache holds the authoritative id<->path index
// (cache.go writes every entry with Expiry == cfg.IDCache.TTL, and for the
// nats-js-kv store the value becomes the bucket-wide MaxAge). Expiring it makes
// the storage provider lose track of existing nodes once entries age out. The
// file-metadata cache, by contrast, is a real cache and must keep its TTL.
func TestDefaultConfigIDCacheHasNoTTL(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.IDCache.TTL != 0 {
		t.Errorf("IDCache.TTL = %s, want 0 (the id<->path index must not expire)", cfg.IDCache.TTL)
	}

	if want := 24 * 60 * time.Second; cfg.FilemetadataCache.TTL != want {
		t.Errorf("FilemetadataCache.TTL = %s, want %s (metadata cache TTL must remain)", cfg.FilemetadataCache.TTL, want)
	}
}
