package ldap

import (
	"crypto/tls"
	"time"
)

// Config holds the basic configuration of the LDAP Connection
type Config struct {
	URI          string
	BindDN       string
	BindPassword string
	TLSConfig    *tls.Config

	RetryMaxCount  int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration


	// PoolSize caps the number of concurrently open connections in the pool. Only used by
	// NewLDAPPool; NewLDAPWithReconnect ignores it. Defaults to defaultPoolSize (5) when <= 0.
	PoolSize int
	// PoolCheckoutTimeout bounds how long a checkout blocks waiting for a connection once the pool
	// is at PoolSize. Only used by NewLDAPPool; NewLDAPWithReconnect ignores it. Defaults to
	// defaultPoolCheckoutTimeout (30s) when <= 0.
	PoolCheckoutTimeout time.Duration
}
