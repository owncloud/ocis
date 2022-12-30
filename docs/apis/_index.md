---
title: APIs
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

{{< toc-tree >}}

Infinite Scale provides a large set of different **application programming interfaces (APIs)**. Infinite Scale is built by microservices. That means many calls to "functions" in the code are remote calls.

Basically we have two different API "universes": [HTTP](http) and [gRPC](grpc_apis).

{{< columns >}} <!-- begin columns block -->

{{< figure src="/ocis/static/http-logo.png" width="70%" alt="Image sourced from https://commons.wikimedia.org/ free of licenses" >}}
<--->

{{< figure src="/ocis/static/grpc-logo.png" width="70%" alt="Image sourced from https://grpc.io/ under CC 4.0 BY license" >}}

{{< /columns >}}


For inter-service-communication we are using mostly gRPC calls because it has some advantages. In the future, clients may decide to use gRPC directly to make use of these advantages.

{{< figure src="/ocis/static/ocis-apis.drawio.svg" class="page-image">}}

## [HTTP](http)

HTTP APIs are mostly used for client <-> server communication. Modern applications are embracing a [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer) software architecture style. REST APIs are using the HTTP protocol to transfer data between clients and servers. All our clients talk to the Server using HTTP APIs. This has legacy reasons and is well-supported across many platforms and technologies. Infinite Scale uses an [HTTP API gateway](../services/proxy) to route client requests to the correct service.

### OpenAPI

It is best practise to define APIs and their behavior by a spec. We are using the OpenAPI standard for all new APIs. The [OpenAPI Specification](https://swagger.io/specification/), previously known as the Swagger Specification, is a specification for a machine-readable interface definition language for describing, producing, consuming and visualizing RESTful web services. Previously part of the Swagger framework, it became a separate project in 2016, overseen by the OpenAPI Initiative, an open-source collaboration project of the Linux Foundation. Swagger and some other tools can generate code, documentation and test cases from interface files.

### RFC

Some APIs have become a de facto standard and are additionally covered by an [RFC](https://en.wikipedia.org/wiki/Request_for_Comments).

## [gRPC](grpc_apis)

In gRPC, a client application can directly call methods on a server application on a different machine as if it was a local object. This makes it easier to create distributed applications based on microservices. In gRPC we can define a service and specify the methods that can be called remotely. A gRPC client has a stub that provides the same methods and types as the server.
Infinite Scale uses a [gRPC API Gateway](../services/gateway) to route the requests to the correct service.

### Protobuf

gRPC APIs are typically defined by [Protocol buffers](https://developers.google.com/protocol-buffers/docs/overview). The different client and server stubs are created from ``*.proto`` files by code generation tools.

## Versioning

There are different standards for API versioning: Through URL, through request parameter, through custom header and through content negotiation. Ocis uses the versioning by URL concept although this creates a big code footprint. The versioning follows [SemVer](https://semver.org). We update the major version number when breaking changes are needed. Clients can decide which major version they use through the request URL. The specific implementation is documented on each API.

