Enhancement: Introduce service registry cache

We've improved the service registry / service discovery by
setting up registry caching (TTL 20s), so that not every requests
has to do a lookup on the registry.

https://github.com/owncloud/ocis/pull/3833
