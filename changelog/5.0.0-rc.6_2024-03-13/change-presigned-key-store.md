Change: change the default store for presigned keys to nats-js-kv

We wrapped the store service in a micro store implementation and changed the default to the built-in NATS instance.

https://github.com/owncloud/ocis/pull/8419
