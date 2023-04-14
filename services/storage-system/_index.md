---
title: Storage-System
date: 2023-04-14T07:43:48.085294925Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/storage-system
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

Purpose and description to be added

## Table of Contents

* [Caching](#caching)
* [Example Yaml Config](#example-yaml-config)

## Caching

The `storage-system` service caches file metadata via the configured store in `STORAGE_SYSTEM_CACHE_STORE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `redis`: Stores metadata in a configured Redis cluster.
  -   `redis-sentinel`: Stores metadata in a configured Redis Sentinel cluster.
  -   `etcd`: Stores metadata in a configured etcd cluster.
  -   `nats-js`: Stores metadata using the key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in productive enviroments.
1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`, `redis` and `redis-sentinel`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  The `storage-system` service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
4.  When using `redis-sentinel`, the Redis master to use is configured via `STORAGE_SYSTEM_CACHE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.

## Example Yaml Config

{{< include file="services/_includes/storage-system-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/storage-system_configvars.md" >}}

