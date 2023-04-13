---
title: Gateway
date: 2023-04-13T15:21:34.15950507Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/gateway
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The gateway service is an ...

## Table of Contents

* [Caching](#caching)
* [Example Yaml Config](#example-yaml-config)

## Caching

The `gateway` service can use a configured store via `GATEWAY_CACHE_STORE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `redis`: Stores data in a configured redis cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  The gateway service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.

## Example Yaml Config

{{< include file="services/_includes/gateway-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/gateway_configvars.md" >}}

