---
title: "Storage"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Overview

The storage extension wraps [reva](https://github.com/cs3org/reva/) and adds an opinionated configuration to provide two core services for the oCIS platform:
1. A [*Spaces Registry*]({{< ref "./spacesregistry.md" >}}) that acts as a dictionary for storage *Spaces* and their metadata
2. A [*Spaces Provider*]({{< ref "./spacesprovider.md" >}}) that organizes *Resources* in storage *Spaces* and persists them in an underlying *Storage System*

*Clients* will use the *Spaces Registry* to poll or get notified about changes in all *Spaces* a user has access to. Every *Space* has a dedicated `/dav/spaces/<spaceid>` WebDAV endpoint that is served by a *Spaces Provider* which uses a specific reva storage driver to wrap an underlying *Storage System*.

{{< svg src="extensions/storage/static/overview.drawio.svg" >}}

The dashed lines in the diagram indicate requests that are made to authenticate requests or lookup the storage provider:
1. After authenticating a request, the proxy may either use the CS3 `userprovider` or the accounts service to fetch the user information that will be minted into the `x-access-token`.
2. The gateway will verify the JWT signature of the `x-access-token` or try to authenticate the request itself, e.g. using a public link token.

{{< hint warning >}}
The bottom part is lighter because we will deprecate it in favor of using only the CS3 user and group providers after moving some account functionality into reva and glauth. The metadata storage is not registered in the reva gateway to separate metadata necessary for running the service from data that is being served directly.
{{< /hint >}}

## Endpoints and references

In order to reason about the request flow, two aspects in the architecture need to be understood well:
1. What kind of [*namespaces*]({{< ref "./namespaces.md" >}}) are presented at the different WebDAV and CS3 endpoints?
2. What kind of [*resource*]({{< ref "./terminology.md#resources" >}}) [*references*]({{< ref "./terminology.md#references" >}}) are exposed or required: path or id based?
{{< svg src="extensions/storage/static/storage.drawio.svg" >}}
