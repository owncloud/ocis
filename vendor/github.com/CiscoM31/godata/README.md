[![Go](https://github.com/CiscoM31/godata/actions/workflows/go.yml/badge.svg)](https://github.com/CiscoM31/godata/actions/workflows/go.yml)
[![golangci-lint](https://github.com/CiscoM31/godata/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/CiscoM31/godata/actions/workflows/golangci-lint.yml)

GoData
======

This is an implementation of OData in Go. It is capable of parsing an OData
request, and exposing it in a standard way so that any provider can consume
OData requests and produce a response. Providers can be written for general
usage like producing SQL statements for a databases, or very specific uses like
connecting to another API.

Most OData server frameworks are C#/.NET or Java. These require using the CLR or
JVM, and are overkill for a lot of use cases. By using Go we aim to provide a
lightweight, fast, and concurrent OData service. By exposing a generic interface
to an OData request, we hope to enable any backend to expose itself with
an OData API with as little effort as possible.

Status
======

This project is not finished yet, and cannot be used in its current state.
Progress is underway to make it usable, and eventually fully compatible with the
OData V4 specification.

Work in Progress
================

* ~~Parse OData URLs~~
* Create provider interface for GET requests
* Parse OData POST and PATCH requests
* Create provider interface for POST and PATCH requests
* Parse OData DELETE requests
* Create provider interface for PATCH requests
* Allow injecting middleware into the request pipeline to enable such features
  as caching, authentication, telemetry, etc.
* Work on fully supporting the OData specification with unit tests

Feel free to contribute with any of these tasks.

High Level Architecture
=======================

If you're interesting in helping out, here is a quick introduction to the
code to help you understand the process. The code works something like this:

1. A provider is initialized that defines the object model (i.e., metadata), of
   the OData service. (See the example directory.)
2. An HTTP request is received by the request handler in service.go
3. The URL is parsed into a data structure defined in request_model.go
4. The request model is semanticized, so each piece of the request is associated
   with an entity/type/collection/etc. in the provider object model.
5. The correct method and type of request (entity, collection, $metadata, $ref, 
   property, etc.) is determined from the semantic information.
6. The request is then delegated to the appropriate method of a GoDataProvider,
   which will produce a response based on the semantic information, and
   package it into a response defined in response_model.go.
7. The response is converted to JSON and sent back to the client.
