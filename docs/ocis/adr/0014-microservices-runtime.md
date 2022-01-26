---
title: "14. Microservices Runtime"
weight: 14
date: 2022-01-21T12:56:53+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0014-microservices-runtime.md
---

* Status: proposed
* Deciders: @butonic, @micbar, @dragotin, @mstingl, @pmaier1, @fschade
* Date: 2022-01-21

## Context and Problem Statement

In an environment where shipping a single binary beats the deployment wars, framing a whole family of microservices within a package and running it using Go language mechanisms has plenty of value. In such environment, a runtime is necessary to orchestrate the services that run within it.

## Decision Drivers

- Start oCIS microservices with a single command (`ocis server`)
- Clear separation of concerns

## Considered Options

1.The use of frameworks such as:
  - asim/go-micro
  - go-kit/kit
2. Build and synchronize all services in-house.
3. A hybrid solution between framework and in-house.

## Options

### go-kit/kit

Pros
- Large community behind
- The creator is a maintainer of Go, so the code quality is quite high.

Cons
- Too verbose. Ultimately too slow to make progress.
- Implementing a service would require defining interfaces and a lot of boilerplate.

### asim/go-micro

Pros
- Implementation based in swappable interfaces.
- Multiple implementations, either in-memory or through external services
- Production ready

## Decision Outcome

Number 3: A hybrid solution between framework and in-house.

### Design

First of, every ocis service IS a go-micro service, and because go-micro makes use of urfave/cli, a service can be conveniently wrapped inside a subcommand. Writing a supervisor is then a choice. We do use a supervisor to ensure long-running processes and embrace the "let it crash" mentality. The piece we use for this end is called [Suture](https://github.com/thejerf/suture).

The code regarding the runtime can be found pretty isolated [here](https://github.com/owncloud/ocis/blob/d6adb7bee83b58aa3524951ed55872a5f3105568/ocis/pkg/runtime/service/service.go). The runtime itself runs as a service. This is done so messages can be sent to it using the oCIS single binary to control the lifecycle of its services.
