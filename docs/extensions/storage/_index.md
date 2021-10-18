---
title: "Storage"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

This service provides an oCIS extension that wraps [reva](https://github.com/cs3org/reva/) and adds an opinionated configuration to it.

## Architecture Overview

The below diagram shows the oCIS services and the contained reva services within as dashed boxes. In general:
1. A request comes in at the proxy and is authenticated using OIDC.
2. It is forwarded to the oCIS frontend which handles ocs and ocdav requests by talking to the reva gateway using the CS3 API.
3. The gateway acts as a facade to the actual CS3 services: storage providers, user providers, group providers and sharing providers.

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
