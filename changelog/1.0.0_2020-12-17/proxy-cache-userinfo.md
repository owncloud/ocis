Enhancement: Cache userinfo in proxy

Tags: proxy

We introduced caching for the userinfo response. The token expiration is used for cache invalidation if available. Otherwise we fall back to a preconfigured TTL (default 10 seconds).

https://github.com/owncloud/ocis/pull/877
