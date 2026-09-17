Enhancement: Retry LDAP operations against a lagging replica

The graph LDAP backend can now retry operations against a replicated directory
where a write to the primary is followed by a read that lands on a replica which
has not yet caught up. This covers the read-back that recovers a directory-assigned
ID after a create (GRAPH_LDAP_SERVER_UUID enabled), which is retried until the entry
becomes visible. The retry count and backoff are tunable through
`GRAPH_LDAP_RETRY_MAX_COUNT`, `GRAPH_LDAP_RETRY_BASE_DELAY` and
`GRAPH_LDAP_RETRY_MAX_DELAY`; the defaults keep the previous behaviour (a single
immediate retry with no delay), so existing deployments are unaffected.

Retries now distinguish reads from writes: a write is no longer retried on a network
error, which can surface after the request was already sent and would otherwise apply
the mutation twice.

https://github.com/owncloud/ocis/pull/12672
