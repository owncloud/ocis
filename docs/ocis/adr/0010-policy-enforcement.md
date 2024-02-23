---
title: "10. Extension Policies"
weight: 10
date: 2021-06-30T14:00:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0010-policy-enforcement.md
---

* Status: proposed
* Deciders: [@butonic](https://github.com/butonic), [@micbar](https://github.com/micbar), [@dragotin](https://github.com/dragotin), [@hodyroff](https://github.com/hodyroff), [@pmaier1](https://github.com/pmaier1), [@fschade](https://github.com/fschade)
* Date: 2021-06-30

## Context and Problem Statement

There should be a way to impose certain limitations in areas of the code that require licensing. This document researches an approach to achieve this goal, while limiting the scope to the enforcement side of it. The architecture for a policy system must be composed of 2 parts:

1. License creation and validation
2. Enforcement

It is desirable to keep both systems isolated, since the implementation of the latter has to be done within the constraints of the codebase. The alternative is running an enforcement service and have each and every single request evaluating whether the request is valid or not.

## Decision Drivers

- As a team, we want to have the licensing code concentrated in a central module
- We don't want to stop/start the extension whenever a policy is updated (hot reload). It must happen during runtime.

## Considered Options

1. Build the evaluation engine in-house.
2. Use third party libraries such as Open Policy Agent (a CNCF approved project written in Go)

## Decision Outcome

Chosen option: option 2; Use third party libraries such as Open Policy Agent (a CNCF approved project written in Go)

### Positive Consequences

- OPA is production battle tested.
- Built around performance - policies evaluations are no longer than 1ms per request.
- Middleware friendly: we use gRPC clients all over our ecosystem; wrappers (or middlewares) is a viable way to solve this problem instead of a dedicated service or its own package.
- Community support.
- Kubernetes friendly.
- Supports envoy, kong, terraform, traefik, php, node and many more.

### Negative Consequences

- More vendor code inside the binary (larger attack surface, larger footprint [to be quantified] )

## Chosen option approach

Make use of [overloading Open Policy Agent's input](https://www.openpolicyagent.org/docs/latest/external-data/#option-2-overload-input) along with an external storage source (instead of an OPA service) in conjunction with go-micro's gRPC client wrappers (a.k.a. middlewares) to leverage policy rules evaluation.

### Terminology

New terms are defined to refer to new mental models:

- Policy: self-imposed limitation of a piece of software. i.e: "after 20 users limit the use of thumbnails".
- Checkers: in the context of a middleware, a checker is in charge of defining logical conditions that prevent requests (users) from doing an action.
- Policy file: a [rego file](https://www.openpolicyagent.org/docs/latest/policy-language/).
- Policy evaluation: the act of piecing together input (from a request), data (from an external storage) and policies in order to make a decision.

#### Temporary new Interfaces part of the PoC

- IStorage: provides means of extracting data from an external source (in case of the POC an etcd storage cluster).

### External data storages

However, for this to be usable it needs state. The Rego engine works with input and data, where data is essentially a database the input is tried against, in order to expand this poc to include functionality such as counters (i.e: give access to the thumbnails only to 50 users) we need an external storage, and consequentially, Rego needs to have an option to load data from an external storage. There is an entire chapter in the documentation regarding external data: https://www.openpolicyagent.org/docs/latest/external-data/. The most "natural" option (option 5) states:

> OPA includes functionality for reaching out to external servers during evaluation. This functionality handles those cases where there is too much data to synchronize into OPA, JWTs are ineffective, or policy requires information that must be as up-to-date as possible.

This is a natural option because it requires service-to-service communication, and by definition using microservices it should come "natural to us". Another approach is using JWT (which we already use) to encode the necessary data into the JWT and handing it over to rego as "data". The issue with this approach is that depending on the features of the licenses the JWT might grow and be filled with noise and redundancy (this is, unless a new token is issued for licensing purposes).

### Future ideas

[This proof of concept](https://github.com/owncloud/ocis/pull/2236) is very rigid in the sense that the `IStorage` interface only has one implementation that ties it to etcd, meaning running an oCIS cluster without an etcd service will result in a crash. This is by far ideal and less coupled implementations should be done. There is the case of using the storage metadata as a source to store data necessary to the policies, or even using the go-micro store as a kv store to achieve the exact same, since it already runs as its own service. The implementation of this is trivial and left out of the POC since it requires more time than the allotted for this task.

#### Message Broker

This problem perfectly encompasses the use of a message broker, where services such as OCS will emit messages to a bus and only listeners react to them. In this case the following applies:

![message broker](https://i.imgur.com/sa1pANQ.jpg)

The necessary interfaces are provided to us by go-micro, only implementations are to be done.
