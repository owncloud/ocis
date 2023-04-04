---
title: Graph service
date: 2023-04-04T10:27:09.204824053Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/graph
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The graph service provides the Graph API which is a RESTful web API used to access Infinite Scale resources. It is inspired by the [Microsoft Graph API](https://learn.microsoft.com/en-us/graph/use-the-api) and can be used by clients or other services or extensions.

## Table of Contents

* [Manual Filters](#manual-filters)
* [Sequence Diagram](#sequence-diagram)
* [Caching](#caching)
* [Example Yaml Config](#example-yaml-config)

## Manual Filters

Using the API, you can manually filter like for users. See the [Libre Graph API](https://owncloud.dev/libre-graph-api/#/users/ListUsers) for examples in the [developer documentation](https://owncloud.dev). Note that you can use `and` and `or` to refine results.

## Sequence Diagram

The following image gives an overview of the scenario when a client requests to list available spaces the user has access to. To do so, the client is directed with his request automatically via the proxy service to the graph service.
<!-- referencing: https://github.com/owncloud/ocis/pull/3816 ([docs-only] add client protocol overview) -->
<img src="https://github.com/owncloud/ocis/blob/master/services/graph/images/mermaid-graph.svg" width="500" />

## Caching

The `graph` service can use a configured store via `GRAPH_STORE_TYPE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `ocmem`: Advanced in-memory store allowing max size.
  -   `redis`: Stores data in a configured redis cluster.
  -   `redis-sentinel`: Stores data in a configured redis sentinel cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
  -   `nats-js`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in productive enviroments.
1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  The graph service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
4.  When using `redis-sentinel`, the Redis master to use is configured via `GRAPH_CACHE_STORE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.

## Example Yaml Config

{{< include file="services/_includes/graph-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/graph_configvars.md" >}}

