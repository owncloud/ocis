Bugfix: Fix etcd address configuration

The etcd server address in `MICRO_REGISTRY_ADDRESS` was not picked up when etcd was set as service discovery registry `MICRO_REGISTRY=etcd`. Therefore etcd was only working if available on localhost / 127.0.0.1.

https://github.com/owncloud/ocis/pull/1546