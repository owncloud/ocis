package roles

import (
	"time"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
)

// cache is a cache implementation for roles, keyed by roleIDs.
type cache struct {
	sc   sync.Cache
	ttl  time.Duration
}

// newCache returns a new instance of Cache.
func newCache(capacity int, ttl time.Duration) cache {
	return cache{
		ttl:     ttl,
		sc: sync.NewCache(capacity),
	}
}

// get gets a role-bundle by a given `roleID`.
func (c *cache) get(roleID string) *settings.Bundle {
	if ce := c.sc.Load(roleID); ce != nil {
		return ce.V.(*settings.Bundle)
	}

	return nil
}

// set sets a roleID / role-bundle.
func (c *cache) set(roleID string, value *settings.Bundle) {
	c.sc.Store(roleID, value, time.Now().Add(c.ttl))
}