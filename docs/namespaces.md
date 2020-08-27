---
title: "Paths and namespaces"
date: 2020-08-27T11:08:00+01:00
weight: 42
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: namespaces.md
---

{{< toc >}}

When a client makes a request to OCIS the path encoded in the request needs to undergo several transformations before actually hitting the storage.
Discussing the related code and request flow is unnecessarily hard without a clear terminology.
Clarifying the different endpoints, involved services and path transformations creates a common vocabulary. This page serves as an  example and forms the glossary for paths, namespaces and related concepts.

## Endpoints

OCIS implements the two WebDAV endpoints that are present in ownCloud 10:
- `/dav/files/<userid>/` which presents the tree of the user with the given user id.
- `/webdav/` which presents the tree of the currently logged in user. It [does not allow referencing another users files](https://github.com/owncloud/core/issues/12504#issuecomment-65259957) and only exists for [backwards compatability](https://github.com/owncloud/core/issues/12543).

For backwards compatibility you can add the old `remote.php`, so these also work:
- `/remote.php/dav/files/<userid>/`
- `/remote.php/webdav/`

## Paths

Whatever follows the two endpoints is the *relative path* of a file in a users home storage. These two URLs actually point to the same file for the user with the id `u-u-i-d`:
- `https://cloud.example.org/dav/files/<userid>/path/to/file.md`
- `https://cloud.example.org/webdav/path/to/file.md`

Since `path/to/file.md` is only relative to the users home storage we need to prefix it with a path to the users home storage in the CS3 namespace.

## Service flow

Before we can dig into the namespaces we should have a look at the services that 
For WebDAV requests:

{{< mermaid class="text-center">}}
graph LR
client-- WebDAV -->proxy
proxy-- WebDAV  --> ocdav
ocdav-- CS3 --> reva-gateway
reva-gateway-- CS3 --> storageprovider
{{< /mermaid >}}

GET and PUT for WebDAV use the data providers instead of storage providers. For a resumable upload the [TUS protocol](https://tus.io/) is used.

Up and downloads take a different route so they can be scaled independently:

{{< mermaid class="text-center">}}
graph LR
client-- tus -->proxy
proxy-- tus --> ocdav
ocdav-- tus --> reva-datagateway
reva-datagateway-- tus --> dataprovider
{{< /mermaid >}}


## Namespaces


| scenario                                                                               | endpoint                                                                            | relative path               | relevant ocdav config                                               | cs3 namespace                                                        | storage_provider_config                                                                                                                                                                                                                                            | user layout example                                                     | internal path                                                                                                                               |
|----------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------|-----------------------------|---------------------------------------------------------------------|----------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|
| legacy access                                                                          | `remote.php/webdav/path/to/file.txt`                                                | `path/to/file.txt`          | WEBDAV_NAMESPACE=`/home/`                                           | `/home/path/to/file.txt`                                             | mount point: `/home`<br>driver: `eoshome`<br>id: `1284d238-aa92-42ce-bdc4-0b0000009158` *Note: this is the id of the storage provider mounted at `/eos`* <br>layout: `{{substr 0 1 .Id.OpaqueId}}/{{.Id.OpaqueId}}`<br>eos namespace: `/eos/dockertest/reva/users` | `f/f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c`                                | namespace + layout + relative path<br>`/eos/dockertest/reva/users` + `f/f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c` + `path/to/file.txt`          |
| userid matches the logged in user<br>einstein = `f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c` | `remote.php/dav/files/f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c/path/to/file.txt`        | `path/to/einstein-file.txt` | WEBDAV_NAMESPACE=`/home/`                                           | `/home/path/to/einstein-file.txt`                                    | *same as above*                                                                                                                                                                                                                                                    | *same as above*                                                         | namespace + layout + relative path<br>`/eos/dockertest/reva/users` + `f/f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c` + `path/to/einstein-file.txt` |
| userid of another user<br>marie = `4c510ada-c86b-4815-8820-42cdf82c3d51`               | `remote.php/dav/files/4c510ada-c86b-4815-8820-42cdf82c3d51/path/to/anotherfile.txt` | `path/to/marie-file.txt`    | FILES_NAMESPACE=`/eos/{{substr 0 1 .Id.OpaqueId}}/{{.Id.OpaqueId}}` | `/eos/4/4c510ada-c86b-4815-8820-42cdf82c3d51/path/to/marie-file.txt` | mount point: `/eos`<br>driver: `eos`<br>id: `1284d238-aa92-42ce-bdc4-0b0000009158`<br>layout: not used by the `eos` driver<br>eos namespace: `/eos/dockertest/reva/users`                                                                                          | *not used by the eos driver, the relative path is constructed by ocdav* | namespace + relative path<br>`/eos/dockertest/reva/users` + `4/4c510ada-c86b-4815-8820-42cdf82c3d51/path/to/marie-file.txt`                 |
