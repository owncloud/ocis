Bugfix: Fix the LDAP master-ID filter when the base filter is empty

The function `EnhanceFilterWithMasterID` in `ocis-pkg/ldap/filter.go` adds a
master-ID clause to an LDAP filter with OR. An empty filter means "no
restriction". Before this fix, the function replaced an empty filter with a
master-ID-only filter. This action restricted an unrestricted search to
master-ID users only.

The bug broke the search for a user on another instance in multi-instance
deployments. This search uses the Graph API endpoint `/graph/v1.0/users`.
Multi-instance deployments use this search to find a user to share a resource
with.

Now the function keeps an empty filter unchanged. The function still adds a
master-ID clause to a filter that is not empty.

The `graph`, `auth-basic`, `users`, and `groups` services also use this
function to compute their default LDAP user filter when
`OCIS_MULTI_INSTANCE_MASTER_ID` is set. An empty user filter now means these
services resolve or authenticate every user in the shared LDAP directory,
not only master-ID users. Each of these services now rejects startup when
`OCIS_MULTI_INSTANCE_MASTER_ID` is set and its own user filter is empty.
Set an explicit `OCIS_LDAP_USER_FILTER` for every multi-instance deployment.

https://github.com/owncloud/ocis/pull/13100
