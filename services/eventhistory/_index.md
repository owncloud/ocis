---
title: Eventhistory Service
date: 2023-03-31T15:38:43.885917074Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/eventhistory
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The `eventhistory` consumes all events from the configured event system like NATS, stores them and allows other services to retrieve them via an eventid.

## Table of Contents

* [Prerequisites](#prerequisites)
* [Consuming](#consuming)
* [Storing](#storing)
* [Retrieving](#retrieving)
* [Example Yaml Config](#example-yaml-config)

## Prerequisites

Running the eventhistory service without an event sytem like NATS is not possible.

## Consuming

The `eventhistory` services consumes all events from the configured event sytem.

## Storing

The `eventhistory` service stores each consumed event via the configured store in `EVENTHISTORY_STORE_TYPE`. Possible stores are:
  -   `mem`: Basic in-memory store and the default.
  -   `ocmem`: Advanced in-memory store allowing max size.
  -   `redis`: Stores data in a configured redis cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
  -   `nats-js`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in productive enviroments.
1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  Events stay in the store for 2 weeks by default. Use `EVENTHISTORY_RECORD_EXPIRY` to adjust this value.
4.  The eventhistory service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.

## Retrieving

Other services can call the `eventhistory` service via a grpc call to retrieve events. The request must contain the eventid that should be retrieved.

## Example Yaml Config

{{< include file="services/_includes/eventhistory-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/eventhistory_configvars.md" >}}

