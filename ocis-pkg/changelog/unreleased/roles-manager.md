Change: Roles manager

We combined the roles middleware and cache into a roles manager. The manager doesn't expose the cache anymore and manages
the state of the cache by fetching roles from the role service which don't exist in the cache, yet.

https://github.com/owncloud/ocis-pkg/pull/60
