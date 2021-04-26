---
title: "5. Account Management through CS3 API"
date: 2021-04-12T15:00:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0004-cs3api-user-management.md
---

# 5. Account Management via CS3 API

* Status: proposed
* Deciders: @refs, @butonic, @micbar, @dragotin, @pmaier1
* Date: 2021-04-12

Technical Story: [Write only management API for User and Group resources](https://github.com/cs3org/cs3apis/pull/119)

## Context and Problem Statement

What would be a more effective way of using network resources and handle account management within the oCIS-Reva ecosystem? Separating account management at the service level is pragmatic and allows for fast iterations, but also steadily accumulates inconsistencies and bloats technical debt.

## Decision Drivers

* Reduce number of network calls.
* Reduce number of services (merge Account + GLAuth from ADR-0003).
* Formalize account management at the API level.

## Considered Options

* Account management delegated to vendors.
* Add account management to the CS3 API.

## Decision Outcome

Chosen option: "Add account management to the CS3 API". Making the API declare an API for account management will not only allow a deployment to fail fast (as in: the management node is not running) but would also centralize all management operations that should happen to be constrained within the Reva context. Constrained operations *SHOULD* be by definition more secure, or at least as secure as the rest of the system.

### Positive Consequences

* More resilient API.
  * Because account management is considered a "first class citizen" changes are forced to go through a more exhaustive revision process.
* Removing Accounts from search users<sup>1</sup>.
* Replace the provisioning API in favor of the new Reva Admin node.

(1) the current vendor implementation of searching a user (i.e: when sharing a resource) relies directly on the accounts service, since this is the only source of truth. Searching a user looks like:

```
┌────────────────────────────────────────┐
│user search (no LDAP)                   │
│                                        │
│    ┌──────────┐                        │
│    │          │                        │
│    │  proxy   │                        │
│    │          │        ┌ ─ ─ ─ ─ ─ ┐   │
│    └──────────┘         go-micro       │
│          ▲             │           │   │
│          │                   Λ         │
│          ▼             │    ╱ ╲    │   │
│    ┌──────────┐            ╱   ╲       │
│    │          │        │  ╱     ╲  │   │
│    │   ocs    │◀──(1)───▶registry▏     │
│    │          │        │  ╲     ╱  │   │
│    └──────────┘            ╲   ╱       │
│          ▲             │    ╲ ╱    │   │
│          │                   V         │
│          │             │           │   │
│          │                             │
│          │             └ ─ ─ ─ ─ ─ ┘   │
│          │                             │
│          │                             │
│          │              ┌──────────┐   │
│          │              │          │   │
│          └─────────────▶│ accounts │   │
│                         │          │   │
│                         └──────────┘   │
│                                        │
│                                        │
│(1) ocs requests a connection to the    │
│accounts service to the registry        │
│                                        │
└────────────────────────────────────────┘
```

Whereas, as a result of ADR-0003 and this ADR, we can simplify and improve this design:

```
┌─────────────────────────────────────────────┐
│user search                                  │
│                                             │
│                                             │
│      ┌──────────┐                           │
│      │          │                           │
│      │  proxy   │                           │
│      │          │                           │
│      └──────────┘                           │
│            │                                │
│            ▼                                │
│      ┌──────────┐                           │
│      │          │                           │
│      │   ocs    │                           │
│      │          │                           │
│      └──────────┘                           │
│            │                                │
│            │                                │
│ ┌ ─ ─ ─ ─ ─│─ ─ ─ ─    ┌ ─ ─ ─ ─ ─ ─ ─ ─ ┐  │
│  reva      ▼       │    IDM                 │
│ │    ┌──────────┐      │   ┌──────────┐  │  │
│      │          │  │       │          │     │
│ │    │  users   │◀─────┼──▶│  GLAuth  │  │  │
│      │          │  │       │          │     │
│ │    └──────────┘      │   └──────────┘  │  │
│                    │                        │
│ └ ─ ─ ─ ─ ─ ─ ─ ─ ─    └ ─ ─ ─ ─ ─ ─ ─ ─ ┘  │
│                                             │
└─────────────────────────────────────────────┘
```

And instead rely on the already existing Reva users provider.


## Pros and Cons of the Options

### Account management delegated to vendors

* Good, because it allows for fast iterations.
* Bad, because account management happens outside of the Reva process. This can potentially end up in invalid account creation / deletion / updates.
  * An example with the existing Accounts service is that any client can fire CRUD accounts requests to the Accounts service as long as the client knows where the server is running and provides with an Authorization header (only required by the proxy). This request totally bypasses Reva middlewares and therefore any security measures that should be enforced by the entire system.
* Bad, because leaves teams the task of designing and implementing a way of dealing with account management. Ideally one schema should be provided / suggested.

Creating an account using the first option looks currently is implemented in vendors as:

```
┌──────────────────────────────────────────────────┐
│ creating a user (webui)                          │
│                                                  │
│       ┌──────────┐                               │
│       │          │                               │
│       │  proxy   │                               │
│       │          │                               │
│       └──────────┘                               │
│             │                                    │
│             │                                    │
│  /api/v0/accounts/accounts-create                │
│             │                                    │
│             │                                    │
│             │                                    │
│             ▼                                    │
│       ┌──────────┐                               │
│       │          │                               │
│       │ accounts │                               │
│       │          │                               │
│       └──────────┘                               │
│                                                  │
│ note that while doing CRUD operations changes    │
│ are instantly reflected for the IDP since out of │
│ the box oCIS uses an accounts backend for        │
│ GLAuth.                                          │
└──────────────────────────────────────────────────┘
```

As explained before, during this flow no Reva middlewares are ran. Creating an account will only use the embedded accounts js file alongside a minted jwt token (by the oCIS proxy) to communicate with the accounts service.

### Add account management to the CS3 API

* Good, because it solidifies what the CS3 API can or cannot do, and account management should be handled at the API level since ultimately accounts would contain a mix of required CS3 and vendor-specific attributes.
* Good, because it centralizes account management and constrains it within the Reva boundaries.
* Good, because there is a clear separation of concerns on what is accounts management logic.
* Good, because we already designed [a similar API for the accounts service](https://github.com/owncloud/ocis/blob/master/accounts/pkg/proto/v0/accounts.proto#L42-L85) the only difference being we (vendors) [define their own messages](https://github.com/owncloud/ocis/blob/master/accounts/pkg/proto/v0/accounts.proto#L252-L408).
  * The API would fully include CRUD methods
* Bad, because development cycles are larger.
  * an example flow will be: `update api > run prototool > publish language specific packages > update dependencies to fetch latest version of the package > utilize the new changes`.

The new account management workflow will result in:
```
┌───────────────────────────────────────────────────┐
│creating a user (webui)                            │
│ - maintain the same route for compatibility       │
│                                                   │
│      ┌──────────┐                                 │
│      │          │                                 │
│      │  proxy   │                                 │
│      │          │                                 │
│      └──────────┘                                 │
│            │                                      │
│            │                                      │
│   /api/v0/accounts/accounts-create                │
│            │                                      │
│            │                                      │
│ ┌ ─ ─ ─ ─ ─│─ ─ ─ ─ ─ ─ ─ ─ ┐  ┌ ─ ─ ─ ─ ─ ─ ─ ─  │
│  Reva      │                    IDM             │ │
│ │          │                │  │                  │
│            ▼                                    │ │
│ │    ┌──────────┐           │  │   ┌──────────┐   │
│      │          │                  │          │ │ │
│ │    │  admin   │───────────┼──┼──▶│  GLAuth  │   │
│      │          │                  │          │ │ │
│ │    └──────────┘           │  │   └──────────┘   │
│                                                 │ │
│ │                           │  │                  │
│  ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─    ─ ─ ─ ─ ─ ─ ─ ─ ┘ │
│                                                   │
│                                                   │
│                                                   │
│                                                   │
│                                                   │
│                                                   │
│ an example of a driver could be GLAuth            │
│ implementing the user management portion of the   │
│ GraphAPI                                          │
└───────────────────────────────────────────────────┘
```

This flow allows Reva and oCIS Proxy to run any middleware logic in order to validate a request. The communication between the proposed Admin api (CS3 API messages) and the IDM (GLAuth) are specific to the _drivers_.
