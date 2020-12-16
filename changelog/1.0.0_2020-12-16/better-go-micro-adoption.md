Enhancement: Better adopt Go-Micro

Tags: ocis

There are a few building blocks that we were relying on default behavior, such as `micro.Registry` and the go-micro client. In order for oCIS to work in any environment and not relying in black magic configuration or running daemons we need to be able to:

- Provide with a configurable go-micro registry.
- Use our own go-micro client adjusted to our own needs (i.e: custom timeout, custom dial timeout, custom transport...)

This PR is relying on 2 env variables from Micro: `MICRO_REGISTRY` and `MICRO_REGISTRY_ADDRESS`. The latter does not make sense to provide if the registry is not `etcd`.

The current implementation only accounts for `mdns` and `etcd` registries, defaulting to `mdns` when not explicitly defined to use `etcd`.

https://github.com/owncloud/ocis/pull/840
