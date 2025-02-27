---
title: "15. oCIS Event System"
weight: 15
date: 2022-02-01T12:56:53+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0015-events.md
---

* Status: proposed
* Deciders: [@butonic](https://github.com/butonic), [@micbar](https://github.com/micbar), [@dragotin](https://github.com/dragotin), [@c0rby](https://github.com/c0rby), [@wkloucek](https://github.com/wkloucek)
* Date: 2022-01-21

## Context and Problem Statement

### Overview

To be able to implement simple, flexible and independent inter service communication there is the idea to implement an event system in oCIS. A service can send out events which are received by one or more other services. The receiving service can cause different kinds of actions based on the event by utilizing the information that the event carries.

### Example: Email Notification

A simple example is the notification feature for oCIS: Users should receive an email when another user shares a file with them. The information, that the file was shared should go out as an event from a storage provider or share manager, carrying the information which file was shared to which receiver. A potential notification service that sends out the email listens to these kinds of events and sends the email out once on every received event of that specific type.

## Decision Drivers

* Events are supposed to decouple services and raise flexibility, also considering extensions that are not directly controlled by the ownCloud project.
* Events should bring flexibility in the implementation of sending and receiving services.
* Events should not obsolete other mechanisms to communicate, i.e. grpc calls.
* Sending an event has to be as little resource consuming for the sender as possible.
* Events are never user visible.

## Considered Options

1. Lightweight Events with Event Queue and "At most once" QoS
2. As 1., but with "At least once" QoS

## Options

### 1. Lightweight Events with Event Queue and "At most once" QoS

Reva will get a messaging service that is available to all services within oCIS and Reva. It is considered as one of the mandatory services of the oCIS system. If the messaging backend is not running, neither Reva nor oCIS can be considered healthy and should shut down.

All oCIS- and Reva-services can connect to the messaging bus and send so-called events. The sender gets an immediate return if handing the event to the message bus was successful or not.

The sender can not make any assumptions when the message is delivered to any receiving service. Depending on the QoS model (as proposed as alternatives in this ADR) it might even be not guaranteed that the event is delivered at all. Also, the sender can not know if zero, one or many services are listening to that event.

#### Event Data

Events are identified by their namespace and their respective name. The namespace is delimited by dots and starts with either "reva" or "ocis" or a future extension name. It is followed by the name of the sending service and an unique name of the event.

Example: `ocis.ocdav.delete` - an event with that name sent out if an WebDAV DELETE request arrived in the oCDav service.

An event can carry a payload which is encoded as json object. (See for example [NATS](https://docs.nats.io/using-nats/developer/sending/structure) ). There are no pre-defined members in that object, it is fully up to the sender which data will be included in the payload. Receivers must be robust to deal with changes.

#### Quality of Service

Events are sent with "At most once" quality of service. That means, if a receiver is not present at the moment of publishing it might not receive the event. That requires that the sender and the receiver must have functionality to back up the situation that events were missed. That adds more state to the services because they always need to behave like a [FISM](https://en.wikipedia.org/wiki/Finite-state_machine). Given that the event queue can be considered the backbone of the system, it is unlikely that it is not running.

#### Transactions

The described way of inter service communication with events is not transactional. It is not supposed to be, but only provides a lightweight, loosely coupled way to "inform".

If transactions are required, proper synchronous GRPC API calls should be used. Another way would be to build asynchronous flows with request- and reply events as in [saga pattern](https://microservices.io/patterns/data/saga.html). That is only recommended for special cases.

#### Pros

* Simple setup
* Flexible way of connecting services
* Stateless event queue
* "State of the art" pattern in microservices architectures

#### Cons

* Over engineering: Can we do without an extra message queue component?
* Messages might get lost, so that eventual consistency is endangered
* A service needs to hold more state to ensure consistency
* Message queue needs to be implemented in Reva

### 2. Lightweight Events with Event Queue and "At-least once" QoS

Exactly as described above, but with a higher service level quality.

#### Quality of Service

Events are sent with "At least once" quality of service. That means the events will remain in the queue until they are received by all receivers. This puts more responsibility on the event bus and adds state to the events. Given that the event queue can be considered the backbone of the system, it is required to be running.

#### Pros

* Better service level: Messages do not get lost
* Simplifies the design of the microservices because the events are "fire-and-forget"
* Events would be idempotent. If a service goes down the events will stay in the queue until they are consumed

#### Cons

* Stateful event system with higher cost in terms of compute and storage
* The queue could become a bottleneck and needs to be scaled

## Decision Outcome

### Design
