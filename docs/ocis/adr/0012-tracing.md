---
title: "12. Tracing"
weight: 12
date: 2021-08-17T12:56:53+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0012-tracing.md
---

* Status: proposed
* Deciders: @butonic, @micbar, @dragotin, @mstingl, @pmaier1, @fschade
* Date: 2021-08-17

## Context and Problem Statement

At the time of this writing we are in a situation where our logs have too much verbosity, rendering impossible or rather difficult to debug an instance. For this reason we are giving some care to our traces by updating dependencies from OpenCensus to OpenTelemetry.

## Decision Drivers

- We don't want to rely only on logs to debug an instance.
- Logs are too verbose.
- Since we have micro-services, we want to holistically understand a request.

## Considered Options

- Trim down logs
- Use OpenCensus
- Migrate to OpenTelemetry

## Decision Outcome

Chosen option: option 3; Migrate to OpenTelemetry. OpenCensus is deprecated, and OpenTelemetry is the merger from OpenCensus and OpenTelemetry and the most recent up-to-date spec.

### Positive Consequences

- Fix the current state of the traces on Reva.
- Add more contextual information on a span for a given request.
- Per-request filtering with the `X-Request-Id` header.
- Group the supported tracing backends to support Jaeger only for simplicity.

## Chosen option approach

- A trace is a tree, and the proxy will create the root trace and propagate it downstream.
- The Root trace will log the request headers.
- The unit that ultimately does the work will log the result of the operation if success.
- The unit that ultimately does the work will change the state of the span to error if any occurred.


With this premises, this is by no means a fixed document and the more we learn about the usage of an instance the more context we can add to the traces.
