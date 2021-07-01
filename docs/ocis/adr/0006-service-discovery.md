---
title: "6. Service Discovery within oCIS and Reva"
date: 2021-04-19T13:00:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0006-service-discovery.md
---

* Status: accepted
* Deciders: @refs, @butonic, @micbar, @dragotin, @pmaier1
* Date: 2021-04-19

Technical Story: [Introduce Named Services.](https://github.com/cs3org/reva/pull/1509)

## Context and Problem Statement

Reva relies heavily on config files. A known implication of this approach are having to know a-priori where a service is running (host + port). We want to move away from hardcoded values and rely instead on named services for service discovery. Furthermore, we would like both platforms (Reva + oCIS) to have the same source of truth at any given time, not having one to notify the other whenever a service status changes.

## Decision Drivers

* Avoid a-priori knowledge of services.
* Ease of scalability.
* Always up-to-date knowledge of the running services on a given deployment (a service registry doesn't have to necessarily be running on the same machine / network)

## Considered Options

* Hardcoded tuples of hostname + port
* Dynamic service registration

## Decision Outcome

Chosen option: "Dynamic service registration". There were some drawbacks regarding this due to introducing go-micro to Reva was from start an issue. Given the little usage of go-micro we need, we decided instead to define our very own [Registry interface](https://github.com/refs/reva/blob/58d013a7509d1941834e1bc814e9a9fa8bff00b1/pkg/registry/registry.go#L22-L35) on Reva and extended the runtime arguments to [allow for injecting a registry](https://github.com/refs/reva/blob/58d013a7509d1941834e1bc814e9a9fa8bff00b1/cmd/revad/runtime/option.go#L53-L58).

### Positive Consequences

* Having dynamic service registration delegates the entire lifecycle of finding a process to the service registry.
* Removing a-priori knowledge of hostname + port for services.
* Marrying go-micro's registry and a newly defined registry abstraction on Reva.
* We will embrace go-micro interfaces by defining a third merger interface in order to marry go-micro registry and rega revistry.
* The ability to fetch a service node relying only on its name (i.e: com.owncloud.proxy) and not on a tuple hostname + port that we rely on being preconfigured during runtime.
* Conceptually speaking, a better framework to tie all the services together. Referring to services by names is less overall confusing than having to add a service name + where it is running. A registry is agnostic to "where is it running" because it, by definition, keeps track of this specific question, so when speaking about design or functionality, it will ease communication.

## Pros and Cons of the Options

### Hardcoded tuples of hostname + port

* Good, because firewalls are easier to configure since IP are static.
* Good, because the mental model required is easier to grasp as IP addresses can be easily bundled.
* Bad, because it requires thorough planning of ports.

### Dynamic service registration

* Good, because it abstracts the use of service lookup away to registry logic from the admin or developer.
* Good, because it allows for, through interfaces, registry injection
  * This means we can have a service registry that we extensively use in oCIS and inject its functionality onto Reva.
* Bad, because it's yet another abstraction.
* Bad, because firewalls are harder to configure with dynamic IPs.f
