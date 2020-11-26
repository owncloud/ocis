Enhancement: Cache basic auth account id in proxy

Tags: proxy

The basic auth middleware now caches account ids. The entry cache gets invalidated after 10 Minutes.
This is useful for scenarios where a lot of basic auth requests with the same username and password happens, for example tests. 

https://github.com/owncloud/ocis/pull/958
