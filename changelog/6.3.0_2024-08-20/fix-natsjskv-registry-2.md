Bugfix: Repair nats-js-kv registry

The registry would always send traffic to only one pod. This is now fixed and load should be spread evenly. Also implements watcher method so the cache can use it.
Internally, it can now distinguish services by version and will aggregate all nodes of the same version into a single service, as expected by the registry cache and watcher.

https://github.com/owncloud/ocis/pull/9734
https://github.com/owncloud/ocis/pull/9726
https://github.com/owncloud/ocis/pull/9656
