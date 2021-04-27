---
title: "Namespaces"
date: 2018-05-02T00:00:00+00:00
weight: 17
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: namespaces.md
---

A *namespace* is a set of paths with a common prefix. Depending on the endpoint you are talking to you will encounter a different kind of namespace:
In ownCloud 10 all paths are considered relative to the users home. The CS3 API uses a global namespace and the *storage providers* use a local namespace with paths relative to the storage providers root.

{{< svg src="extensions/storage/static/namespaces.drawio.svg" >}}

The different paths in the namespaces need to be translated while passing *references* from service to service. While the oc10 endpoints all work on paths we internally reference  shared resources by id, so the shares don't break when a file is renamed or moved inside a *storage space*.

| oc10 namespace                                   | CS3 global namespace                          | storage provider | reference |  content |
|--------------------------------------------------|----------------------------------------|---------|-------------------|-----------------|
| `/webdav/path/to/file.ext` `/dav/files/<username>/path/to/file.ext`                       | `/home/path/to/file.ext`               | home | `/<userlayout>/path/to/file.ext` | currently logged in users home |
| `/webdav/Shares/foo` `/dav/files/<username>/Shares/foo` | `/home/Shares/foo`                     | users | id based access | all users, used to access collaborative shares |
| `/dav/public-files/<token>/rel/path/to/file.ext` | `/public/<token>/rel/path/to/file.ext` | public | id based access | publicly shared files, used to access public links |


{{< hint danger >}}
oCIS currently is configured to jail users into the CS3 `/home` namespace in the oc10 endpoints to mimic ownCloud 10. CernBox has been exposing a global namespace on `/webdav` for years already. The ocs service returns urls that are relative to the CS3 global namespace which makes both scenarios work, but only one of them at a time. Which is why the testsuite hiccups when trying to [Allow full paths targets in reva#1605](https://github.com/cs3org/reva/pull/1605).
{{< /hint >}}


{{< hint warning >}}
In the global CS3 namespaces we plan to move `/home/Shares`, which currently lists all mounted shares of the currently logged in user to a dedicated `/shares` namespace. See [Move shares folder out from home directory to a separate mount reva#1584](https://github.com/cs3org/reva/pull/1584).
{{< /hint >}}
