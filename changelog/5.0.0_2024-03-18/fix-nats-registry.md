Bugfix: Fix nats registry

The nats registry would behave badly when configuring `nats-js-kv` via envvar. Reason is the way go-micro initializes.
It took 5 developers to find the issue and the fix so the details cannot be shared here. Just accept that it is working now

https://github.com/owncloud/ocis/pull/8281
