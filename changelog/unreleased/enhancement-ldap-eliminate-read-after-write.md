Enhancement: Eliminate redundant LDAP read-after-write on create and update

The graph LDAP backend no longer re-reads an entry immediately after writing it
just to recover the entry ID for the response. When oCIS generates the ID itself
(GRAPH_LDAP_SERVER_UUID disabled), the create response is now synthesized from the
data already sent to the directory, and update responses are built by folding the
applied modifications onto the entry that was read before the write. This avoids a
round-trip that, against a replicated directory reached through a proxy, could hit
a lagging replica and fail or return stale data.

When the directory assigns the ID (GRAPH_LDAP_SERVER_UUID enabled), creates keep
the existing read-back, since the generated ID cannot otherwise be recovered.

https://github.com/owncloud/ocis/pull/12617
