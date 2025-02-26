Bugfix: Fix nats registry

Using `nats` as service registry did work, but when a service would restart and gets a new ip it couldn't re-register.
We fixed this by using `"put"` register action instead of the default `"create"`

https://github.com/owncloud/ocis/pull/6881
