---
title: "Namespaces"
date: 2018-05-02T00:00:00+00:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: namespaces.md
---

A *namespace* is a set of paths with a common prefix. Depending on the endpoint you are talking to you will encounter a different kind of namespace:
In ownCloud 10 all paths are considered relative to the users home. The CS3 API uses a global namespace and the *storage providers* use a local namespace with paths relative to the storage providers root.

{{< svg src="extensions/storage/static/namespaces.drawio.svg" >}}

The different paths in the namespaces need to be translated while passing [*references*]({{< ref "./terminology.md#references" >}}) from service to service. While the oc10 endpoints all work on paths we internally reference shared resources by id, so the shares don't break when a file is renamed or moved inside a [*storage space*]({{< ref "./terminology.md#storage-spaces" >}}). The following table lists the various namespaces, paths and id based references:

| oc10 namespace                                   | CS3 global namespace                   | storage provider | reference | content |
|--------------------------------------------------|----------------------------------------|------------------|-----------|---------|
| `/webdav/path/to/file.ext` `/dav/files/<username>/path/to/file.ext`                       | `/home/path/to/file.ext` | home | `/<userlayout>/path/to/file.ext` | currently logged in users home |
| `/webdav/Shares/foo` `/dav/files/<username>/Shares/foo` | `/home/Shares/foo`              | users | id based access | all users, used to access collaborative shares |
| `/dav/public-files/<token>/rel/path/to/file.ext` | `/public/<token>/rel/path/to/file.ext` | public | id based access | publicly shared files, used to access public links |


{{< hint danger >}}
oCIS currently is configured to jail users into the CS3 `/home` namespace in the oc10 endpoints to mimic ownCloud 10. CernBox has been exposing a global namespace on `/webdav` for years already. The ocs service returns urls that are relative to the CS3 global namespace which makes both scenarios work, but only one of them at a time. Which is why the testsuite hiccups when trying to [Allow full paths targets in reva#1605](https://github.com/cs3org/reva/pull/1605).
{{< /hint >}}


{{< hint warning >}}
In the global CS3 namespaces we plan to move `/home/Shares`, which currently lists all mounted shares of the currently logged in user to a dedicated `/shares` namespace. See [below]({{< ref "#cs3-namespaces" >}}) and [Move shares folder out from home directory to a separate mount reva#1584](https://github.com/cs3org/reva/pull/1584).
{{< /hint >}}

## ownCloud namespaces

In contrast to the global namespace of CS3, ownCloud always presented a user specific namespace on all endpoints. It will always list the users private files under `/`. Shares can be mounted at an arbitrary location in the users private spaces. See the [webdav]({{< ref "./architecture#webdav" >}}) and [ocs]({{< ref "./architecture#sharing" >}}) sections for more details end examples.

With the spaces concept we are planning to introduce a global namespace to the ownCloud webdav endpoints. This will push the users private space down in the hierarchy: it will move from `/webdav` to `/webdav/home` or `/webdav/users/<username>`. The related [migration stages]({{< ref "../../ocis/migration.md" >}}) are subject to change.

## CS3 global namespaces

The *CS3 global namespace* in oCIS is configured in the [*storage registry*]({{< ref "./terminology.md#storage-registries" >}}). oCIS uses these defaults:

| global namespace | description |
|-|-|
| `/home` | an alias for the currently logged in uses private space |
| `/users/<userlayout>` | user private spaces |
| `/shares` | a virtual listing of share spaces a user has access to |
| `/public/<token>` | a virtual folder listing public shares |
| `/spaces/<spacename>` | *TODO: project or group spaces* |

Technically, the `/home` namespace is not necessary: the [*storage registry*]({{< ref "./terminology.md#storage-registries" >}}) knows the path to a users private space in the `/users` namespace and the gateway can forward the requests to the responsible storage provider.

{{< hint warning >}}
*@jfd: Why don't we use `/home/<userlayout>` instead of `/users/<userlayout>`. Then the paths would be consistent with most unix systems.
{{< /hint >}}

The `/shares` namespace is used to solve two problems:
- To query all shares the current user has access to the *share manager* can be used to list the resource ids. While the shares can then be navigated by resource id, they will return the relative path in the actual [*storage provider*]({{< ref "./terminology.md#storage-providers" >}}), leaking parent folders of the shared resource.
- When accepting a remote share e.g., for OCM the resource does not exist on the local instance. They are made accessible in the global namespace under the `/shares` namespace.

{{< hint warning >}}
*@jfd: Should we split `/shares` into `/collaborations`, `/ocm` and `/links`? We also have `/public` which uses token based authentication. They may have different latencies or polling strategies? Well, I guess we can cache them differently regardless of the mount point.*
{{< /hint >}}