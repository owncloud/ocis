---
title: gRPC
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/grpc_apis/
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

{{< toc-tree >}}

## **R**emote  &nbsp; **P**rocedure  &nbsp; **C**alls

[gRPC](https://grpc.io) is a modern open source high performance Remote Procedure Call (RPC) framework that can run in any environment. It can efficiently connect services in and across data centers with pluggable support for load balancing, tracing, health checking and authentication. It is also applicable in last mile of distributed computing to connect devices, mobile applications and browsers to backend services.

## Advantages of gRPC

{{< columns >}}
### {{< icon "gauge-high" >}} &nbsp; Performance

gRPC uses http/2 by default and is faster than REST. When using protocol buffers for encoding, the information comes on and off the wire much faster than JSON. Latency is an important factor in distributed systems. JSON encoding creates a noticeable factor of latency. For distributed systems and high data loads, gRPC can actually make an important difference. Other than that, gRPC supports multiple calls via the same channel and the connections are bidirectional. A single connection can transmit requests and responses at the same time. gRPC keeps connections open to reuse the same connection again which prevents latency and saves bandwidth.

<--->
### {{< icon "helmet-safety" >}} &nbsp; Robustness

gRPC empowers better relationships between clients and servers. The rules of communication are strictly enforced. That is not the case in REST calls, where the client and the server can send and receive anything they like and hopefully the other end understands what to do with it. In gRPC, to make changes to the communication, both client and server need to change accordingly. This prevents mistakes specially in microservice architectures.
{{< /columns >}}
{{< columns >}}

### {{< icon "magnifying-glass-plus" >}} &nbsp; Debuggability

gRPC requests are re-using the same context and can be tracked or traced across multiple service boundaries.
This helps to identify slow calls and see what is causing delays. It is possible to cancel requests which cancels
them on all involved services.

<--->
### {{< icon "boxes-stacked" >}} &nbsp; Microservices

gRPC has been evolving and has become the best option for communication between microservices because of its unmatched
performance and its polyglot nature. One of the biggest strengths of microservices is the freedom of programming
languages and technologies. By using gRPC we can leverage all the advantages of strictly enforced communication
standards combined with freedom of choice between different programming languages - whichever would fit best.

{{< /columns >}}

{{< hint type=info title="gRPC Advantages" >}}

- http/2
- protocol buffers
- reusable connections
- multi language support
{{< /hint >}}

## CS3 APIs

{{< figure src="/ocis/static/cs3org.png" >}}

The [CS3 APIs](https://github.com/cs3org/cs3apis) connect storages and application providers.

The CS3 APIs follow Google and Uber API design guidelines, specially on error handling and naming convention. You can read more about these
guidelines at https://cloud.google.com/apis/design/ and https://github.com/uber/prototool/blob/dev/style/README.md.

The CS3 APIs use [Protocol Buffers version 3 (proto3)](https://github.com/protocolbuffers/protobuf) as their
Interface Definition Language (IDL) to define the API interface and the structure of the payload messages.
