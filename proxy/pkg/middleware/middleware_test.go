package middleware

import (
	"github.com/owncloud/ocis/proxy/pkg/cache"
)

const (
	// AccountsKey declares the svcKey for the Accounts service.
	AccountsKey = "accounts"
)

var (
	// svcCache caches requests for given services to prevent round trips to the service
	svcCache = cache.NewCache(
		cache.Size(256),
	)
)
