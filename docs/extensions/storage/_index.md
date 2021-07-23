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

{{< svg src="extensions/storage/static/storage.drawio.svg" >}}
