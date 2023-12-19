---
title: Registry
date: 2023-12-19T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info
geekdocFilePath: registry.md
geekdocCollapseSection: true
---

To be able to communicate with each other, services need to register in a common registry.

## Configuration

The type of registry to use can be configured with the `MICRO_REGISTRY` environment variable. Supported values are:

### `memory`

Setting the environment variable to `memory` starts an inmemory registry. This only works with the single binary deployment.

### `nats-js-kv`

Set the environment variable to `nats-js-kv` (or leave it empty) to use a nats-js key value store as registry.
- Note: If not running build-in nats, `MICRO_REGISTRY_ADDRESS` needs to be set to the address of the nats-js cluster. (Same as `OCIS_EVENTS_ENDPOINT`)
- Optional: Use `MICRO_REGISTRY_AUTH_USERNAME` and `MICRO_REGISTRY_AUTH_PASSWORD` to authenticate with the nats cluster.

This is the default.

### `kubernetes`

When deploying in a kubernetes cluster, the kubernetes registry can be used. Additionally the `MICRO_REGISTRY_ADDRESS` environment
variable needs to be set to the url of the kubernetes registry.

### Deprecated registries

These registries are currently working but will be removed in a later version. It is recommended to switch to a supported one.

- `nats`. Uses a registry based on nats streams. Requires `MICRO_REGISTRY_ADDRESS` to bet set.
- `etcd`. Uses an etcd cluster as registry. Requires `MICRO_REGISTRY_ADDRESS` to bet set.
- `consul`. Uses `HashiCorp Consul` as registry. Requires `MICRO_REGISTRY_ADDRESS` to bet set.
- `mdns`.  Uses multicast dns for registration. This type can have unwanted side effects when other devices in the local network use mdns too.

