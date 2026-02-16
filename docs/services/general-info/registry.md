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

{{< toc >}}

## Configuration

The type of registry to use can be configured with the `MICRO_REGISTRY` environment variable. Supported values are:

### `nats-js-kv` (Default)

Set the environment variable to `nats-js-kv` or leave it empty to use a nats-js key value store as registry.

- Note: If not running build-in nats, `MICRO_REGISTRY_ADDRESS` needs to be set to the address of the nats-js cluster, which is the same value as `OCIS_EVENTS_ENDPOINT`.
- Optional: Use `MICRO_REGISTRY_AUTH_USERNAME` and `MICRO_REGISTRY_AUTH_PASSWORD` to authenticate with the nats cluster.

### `kubernetes`

When deploying in a kubernetes cluster, the Kubernetes registry can be used. Additionally, the `MICRO_REGISTRY_ADDRESS` environment variable needs to be set to the url of the Kubernetes registry.

### `memory`

Setting the environment variable to `memory` starts an in-memory registry. This only works with the single binary deployment.

### Deprecated Registries

These registries are currently working but will be removed in a later version. It is recommended to switch to a supported one.

- `nats`\
  Uses a registry based on nats streams. Requires `MICRO_REGISTRY_ADDRESS` to be set.
- `etcd`\
  Uses an etcd cluster as the registry. Requires `MICRO_REGISTRY_ADDRESS` to be set.
- `consul`\
  Uses `HashiCorp Consul` as registry. Requires `MICRO_REGISTRY_ADDRESS` to be set.
- `mdns`\
  Uses multicast dns for registration. This type can have unwanted side effects when other devices in the local network use mdns too.

