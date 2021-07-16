---
title: "10. Extension Policies"
date: 2021-06-30T14:00:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0010-policy-enforcement.md
---

* Status: proposed
* Deciders: @butonic, @micbar, @dragotin, @hodyroff, @pmaier1, @fschade
* Date: 2021-06-30

## Context and Problem Statement

There should be a way to impose certain limitations in areas of the code that require licensing. This document researches an approach to achieve it.

## Decision Drivers

- as a team, we want to have the licensing code concentrated in a central module
- we don't want to stop/start the extension whenever a policy is updated (hot reload)

## Considered Options

1. Build the evaluation engine in-house.
2. Use third party libraries such as Open Policy Agent (a CNCF aproved project written in Go)

## Decision Outcome

Chosen option: option 2; Use third party libraries such as Open Policy Agent (a CNCF aproved project written in Go)

### Positive Consequences

- OPA is production battle tested.
- built around performance - policies evaluations are no longer than 1ms per request.
- middleware friendly: we use gRPC clients all over our ecosystem; wrappers (or middlewares) is a viable way to solve this problem instead of a dedicated service or its own package.
- community support.
- kubernetes friendly.
- supports envoy, kong, terraform, traefik, php, node and many more.

### Negative Consequences

- more vendor code inside the binary (larger attack surface, larger footprint [to be quantified] )

## Chosen option approach

Make use of [overloading Open Policy Agent's input](https://www.openpolicyagent.org/docs/latest/external-data/#option-2-overload-input) along with an external storage source (instead of an OPA service) in conjunction with go-micro's gRPC client wrappers (a.k.a middlewares) to leverage policy rules evaluation.

### Terminology

New terms are defined to refer to new mental models:

- policy: self-imposed limitation of a piece of software. i.e: "after 20 users limit the use of thumbnails".
- checkers: in the context of a middleware, a checker is in charge of defining logical conditions that prevent requests (users) from doing an action.
- policy file: a [rego file](https://www.openpolicyagent.org/docs/latest/policy-language/).
- policy evaluation: the act of piecing together input (from a request), data (from an external storage) and policies in order to make a decision.

#### Temporary new Interfaces part of the PoC

- IStorage: provides means of extracting data from an external source (in case of the POC an etcd storage cluster).

### Future ideas

[This proof of concept](https://github.com/owncloud/ocis/pull/2236) is very rigid in the sense that the `IStorage` interface only has one implementation that ties it to etcd, meaning running an oCIS cluster without an etcd service will result in a crash. This is by far ideal and less coupled implementations should be done. There is the case of using the storage metadata as a source to store data necessary to the policies, or even using the go-micro store as a kv store to achieve the exact same, since it already runs as its own service. The implementation of this is trivial and left out of the POC since it requires more time than the allotted for this task.

#### Message Broker

This problem perfectly encompasses the use of a message broker, where services such as OCS will emit messages to a bus and only listeners react to them. In this case the following applies:

![message broker](https://i.imgur.com/sa1pANQ.jpg)

The necessary interfaces are provided to us by go-micro, only implementations are to be done.
