package roles

import (
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/sync"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
)

// cache is a cache implementation for roles, keyed by roleIDs.
type cache struct {
	sc  sync.Cache
	ttl time.Duration
}

// newCache returns a new instance of Cache.
func newCache(capacity int, ttl time.Duration) cache {
	return cache{
		ttl: ttl,
		sc:  sync.NewCache(capacity),
	}
}

// get gets a role-bundle by a given `roleID`.
func (c *cache) get(roleID string) *settingsmsg.Bundle {
	if ce := c.sc.Load(roleID); ce != nil {
		return ce.V.(*settingsmsg.Bundle)
	}

	return nil
}

// set sets a roleID / role-bundle.
func (c *cache) set(roleID string, value *settingsmsg.Bundle) {
	c.sc.Store(roleID, value, time.Now().Add(c.ttl))
}
