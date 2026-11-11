Bugfix: Fix the role assignment of multi-instance guests

The proxy now rejects an unknown multi-instance guest role. Before this
fix, the proxy assigned an empty role id when
`OCIS_MULTI_INSTANCE_GUEST_ROLE` named a role that does not exist.

The proxy now trims the spaces in the pairs of
`OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING`. The proxy does not start when the
mapping is set and `OCIS_MULTI_INSTANCE_GUEST_ROLE` is empty.

When a guest role is not in the role cache, the proxy reads the roles again.
Thus a new role works before the cache time ends. The proxy reads the roles
a maximum of one time in 30 seconds.

The proxy writes a debug log line when the roles claim of a guest is missing.
The description of `OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING` now tells that
a claim value is a regular expression.

https://github.com/owncloud/ocis/pull/13250
