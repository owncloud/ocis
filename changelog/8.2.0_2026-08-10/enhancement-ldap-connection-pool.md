Enhancement: Add an opt-in bounded LDAP connection pool

The auth-basic, users, groups and graph services can now be switched from a
single long-lived reconnecting LDAP connection to a bounded pool of connections,
so concurrent requests no longer serialize on one socket. Connections are dialed
and bound lazily on checkout, unhealthy connections are discarded and lazily
re-dialed rather than eagerly reconnected, and checkout blocks with a configurable
timeout once the pool is exhausted.

Pooling is off by default and fully backwards compatible. Enable it per service
via 'OCIS_LDAP_POOL_ENABLED' (or the service-specific '<SERVICE>_LDAP_POOL_ENABLED'
override), and tune it with 'OCIS_LDAP_POOL_SIZE' (default 5) and
'OCIS_LDAP_POOL_CHECKOUT_TIMEOUT' (default 30s).

The graph service's identity backend now shares the same LDAP client
implementation used by the reva auth/user/group managers instead of maintaining
its own separate reconnecting client.

https://github.com/owncloud/ocis/pull/12688
