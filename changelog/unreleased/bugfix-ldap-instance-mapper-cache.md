Bugfix: Cache LDAP instance mapper lookups

In multi-instance deployments, resolving a user's instance name/ID during
`GET /graph/v1.0/users` (and group member expansion) issued a fresh, uncached
LDAP search per instance/guest attribute value on every request. Under load
this multiplied into large numbers of redundant LDAP round-trips per page of
users, saturating the LDAP connection pool and causing request timeouts. The
LDAP identity backend now caches instance mapper lookups, including negative
(not-found) results, for a configurable TTL
(`OCIS_LDAP_INSTANCE_MAPPER_CACHE_TTL`, default 60s).
