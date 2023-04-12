---
title: OCS Service
date: 2023-04-12T03:38:13.90523951Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/ocs
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The ocs service is an ...

## Table of Contents

* [Caching](#caching)
* [Example Yaml Config](#example-yaml-config)

## Caching

The `ocs` service can use a configured store via `OCS_STORE_TYPE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `ocmem`: Advanced in-memory store allowing max size.
  -   `redis`: Stores data in a configured redis cluster.
  -   `redis-sentinel`: Stores data in a configured redis sentinel cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
  -   `nats-js`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in productive enviroments.
1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  The ocs service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
4.  When using `redis-sentinel`, the Redis master to use is configured via `OCS_CACHE_STORE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.

## Example Yaml Config

{{< include file="services/_includes/ocs-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/ocs_configvars.md" >}}

