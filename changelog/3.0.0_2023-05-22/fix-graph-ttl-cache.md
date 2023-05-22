Bugfix: Reduced default TTL of user and group caches in graph API

We reduced the default TTL of the cache for user and group information on the
/drives endpoints to 60 seconds. This fixes in issue where outdated information
was show on the spaces list for a very long time.

https://github.com/owncloud/ocis/issues/6320
