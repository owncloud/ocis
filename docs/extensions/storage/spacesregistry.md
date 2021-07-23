---
title: "Spaces Registry"
date: 2018-05-02T00:00:00+00:00
weight: 9
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: spacesregistry.md
---

{{< hint warning >}}

The current implementation in oCIS might not yet fully reflect this concept. Feel free to add links to ADRs, PRs and Issues in short warning boxes like this.

{{< /hint >}}

## Storage Space Registries

A storage *spaces registry* manages the [*namespace*]({{< ref "./namespaces.md" >}}) for a *user*: it is used by *clients* to look up storage spaces a user has access to, the `/dav/spaces` endpoint to access it via WabDAV, and where the client should mount it in the users personal namespace.

{{< svg src="extensions/storage/static/spacesregistry.drawio.svg" >}}

