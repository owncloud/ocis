Bugfix: Fix the LDAP master-ID filter when the base filter is empty

The function EnhanceFilterWithMasterID adds a master-ID clause to an LDAP
filter with OR. An empty filter means "no restriction". Before this fix, the
function replaced an empty filter with a master-ID-only filter. This action
restricted an unrestricted search to master-ID users only.

The bug broke the search for a user on another instance in multi-instance
deployments. This search uses the Graph API endpoint `/graph/v1.0/users`.
Multi-instance deployments use this search to find a user to share a resource
with.

Now the function keeps an empty filter unchanged. The function still adds a
master-ID clause to a filter that is not empty.

https://github.com/owncloud/ocis/pull/13096
