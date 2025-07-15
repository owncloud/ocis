Enhancement: Add glauth fallback backend

We introduced the `fallback-datastore` config option and the corresponding options to allow configuring a simple chain of two handlers.

Simple, because it is intended for bind and single result search queries. Merging large sets of results is currently out of scope. For now, the implementation will only search the fallback backend if the default backend returns an error or the number of results is 0. This is sufficient to allow an IdP to authenticate users from ocis as well as owncloud 10 as described in the [bridge scenario](https://owncloud.github.io/ocis/deployment/bridge/).

https://github.com/owncloud/ocis/pull/649
https://github.com/owncloud/ocis-glauth/issues/18
