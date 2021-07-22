---
title: "Proposed Changes"
date: 2018-05-02T00:00:00+00:00
weight: 18
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: proposedchanges.md
---

Some architectural changes still need to be clarified or changed. Maybe an ADR is in order for all of the below.

## Reva Gateway changes

## A dedicated shares storage provider

Currently, the *gateway* treats `/home/shares` different than any other path: it will stat all children and calculate an etag to allow clients to discover changes in accepted shares. This requires the storage provider to cooperate and provide this special `/shares` folder in the root of a users home when it is accessed as a home storage, which is a config flag that needs to be set for every storage driver.

The `enable_home` flag will cause drivers to jail path based requests into a `<userlayout>` subfolder. In effect it divides a storage provider into multiple [*storage spaces*]({{< ref "#storage-spaces" >}}): when calling `CreateHome` a subfolder following the `<userlayout>` is created and market as the root of a users home. Both, the eos and ocis storage drivers use extended attributes to mark the folder as the end of the size aggregation and tree mtime propagation mechanism. Even setting the quota is possible like that. All this literally is a [*storage space*]({{< ref "#storage-spaces" >}}).

We can implement [ListStorageSpaces](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ListStorageSpacesRequest) by either
- iterating over the root of the storage and treating every folder following the `<userlayout>` as a `home` *storage space*, 
- iterating over the root of the storage and treating every folder following a new `<projectlayout>` as a `project` *storage space*, or
- iterating over the root of the storage and treating every folder following a generic `<layout>` as a *storage space* for a configurable space type, or
- we allow configuring a map of `space type` to `layout` (based on the [CreateStorageSpaceRequest](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.CreateStorageSpaceRequest)) which would allow things like
```
home=/var/lib/ocis/storage/home/{{substr 0 1 .Owner.Username}}/{{.Owner.Username}}
spaces=/spaces/var/lib/ocis/storage/projects/{{.Name}}
```

This would make the `GetHome()` call return the path to the *storage provider* including the relative path to the *storage space*. No need for a *storage provider* mounted at `/home`. This is just a UI alias for `/users/<userlayout>`. Just like a normal `/home/<username>` on a linux machine.

But if we have no `/home` where do we find the shares, and how can clients discover changes in accepted shares?

The `/shares` namespace should be provided by a *shares storage provider* that lists all accepted shares for the current user... but what about copy pasting links from the browser? Well this storage is only really needed to have a path to ocm shares that actually reside on other instances. In the UI the shares would be listed by querying a *share manager*. It returns ResourceIds, which can be stated to fetch a path that is then accessible in the CS3 global namespace. Two caveats:
- This only works for resources that are actually hosted by the current instance. For those it would leak the parent path segments to a shared resource.
- For accepted OCM shares there must be a path in the [*CS3 global namespace*]({{< ref "./namespaces.md#cs3-global-namespaces" >}}) that has to be the same for all users, otherwise they cannot copy and share those URLs.

Work on this is done in https://github.com/cs3org/reva/pull/1846

### The gateway should be responsible for path transformations

Currently, storage providers are aware af their mount point, coupling them tightly with the gateway.

Tracked in https://github.com/cs3org/reva/issues/578

Work is done in https://github.com/cs3org/reva/pull/1866

## URL escaped string representation of a CS3 reference

For the `/dav/spaces/` endpoint we need to encode the *reference* in a url compatible way. 
1. We can separate the path using a `/`: `/dav/spaces/<spaceid>/<path>`
2. The `spaceid` currently is a cs3 resourceid, consisting of `<storageid>` and `<nodeid>`. Since the nodeid might contain `/` eg. for the local driver we have to urlencode the spaceid.

To access resources by id we need to make the `/dav/meta/<resourceid>` able to list directories... Otherwise id based navigation first has to look up the path. Or we use the libregraph api for id based navigation.

A *reference* is a logical concept. It identifies a [*resource*]({{< ref "#resources" >}}) and consists of a `<resource_id>` and a `<path>`. A `<resource_id>` consists of a `<storage_id>` and a `<node_id>`. They can be concatenated using the separators `!` and `:`:
```
<storage_id>!<node_id>:<path>
```
While all components are optional, only three cases are used:
| format | example | description |
|-|-|-|
| `!:<absolute_path>` | `!:/absolute/path/to/file.ext` | absolute path | 
| `<storage_space>!:<relative_path>` | `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!:path/to/file.ext` | path relative to the root of the storage space | 
| `<storage_space>!<root>:<relative_path>` | `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!c3cf23bb-8f47-4719-a150-1d25a1f6fb56:to/file.ext` | path relative to the specified node in the storage space, used to reference resources without disclosing parent paths |

`<storage_space>` should be a UUID to prevent references from breaking when a *user* or [*storage space*]({{< ref "#storage-spaces" >}}) gets renamed. But it can also be derived from a migration of an oc10 instance by concatenating an instance identifier and the numeric storage id from oc10, e.g. `oc10-instance-a$1234`.

A reference will often start as an absolute/global path, e.g. `!:/home/Projects/Foo`. The gateway will look up the storage provider that is responsible for the path

| Name | Description | Who resolves it? |
|------|-------------|-|
| `!:/home/Projects/Foo` | the absolute path a client like davfs will use. | The gateway uses the storage registry to look up the responsible storage provider |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!:/Projects/Foo` | the `storage_space` is the same as the `root`, the path becomes relative to the root | the storage provider can use this reference to identify this resource |

Now, the same file is accessed as a share
| Name | Description |
|------|-------------|
| `!:/users/Einstein/Projects/Foo` | `Foo` is the shared folder |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:` | `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a` is the id of `Foo`, the path is empty |


The `:`, `!` and `$` are chosen from the set of [RFC3986 sub delimiters](https://tools.ietf.org/html/rfc3986#section-2.2) on purpose. They can be used in URLs without having to be encoded. In some cases, a delimiter can be left out if a component is not set:
| reference | interpretation |
|-|-|
| `/absolute/path/to/file.ext` | absolute path, all delimiters omitted |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!path/to/file.ext` | relative path in the given storage space, root delimiter `:` omitted |
| `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:to/file.ext` | relative path in the given root node, storage space delimiter `!` omitted |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:` | node id in the given storage space, `:` must be present |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62` | root of the storage space, all delimiters omitted, can be distinguished by the `/` |

## space providers 
When looking up an id based resource the reference must use a logical space id, not a CS3 resource id. Otherwise id based requests, which only have a resourceid consisting of a storage id and a node id cannot be routed to the correct storage provider if the storage has moved from one storage provider to another. 

if the registry routes based on the storageid AND the nodeid it has to keep a cache of all nodeids in order to route all requests for a storage space (which consists of storage it + nodeid) to the correct storage provider. the correct resourceid for a node in a storage space would be `<storageid>$<rootnodeid>!<nodeid>`. The `<storageid>$<rootnodeid>` part allow the storage registry to route all id based requests to the correct storage provider. This becomes relevant when the storage space was moved from one storage provider to another. The storage space id remains the same, but the internal address and port change.

TODO discuss to clarify further

## Storage drivers

### allow clients to send a uuid on upload
iOS clients can only queue single requests to be executed in the background. They queue an upload and need to be able to identify the uploaded file after it has been uploaded to the server. The disconnected nature of the connection might cause workflows or manual user interaction with the file on the server to move the file to a different place or changing the content while the device is offline. However, on the device users might have marked the file as favorite or added it to other iOS specific collections. To be able to reliably identify the file the client can generate a `uuid` and attach it to the file metadata during the upload. While it is not necessary to look up files by this `uuid` having a second file id that serves exactly the same purpose as the `file id` is redundant.

Another aspect for the `file id` / `uuid` is that it must be a logical identifier that can be set, at least by internal systems. Without a writeable fileid we cannot restore backups or migrate storage spaces from one storage provider to another storage provider.

Technically, this means that every storage driver needs to have a map of a `uuid` to an internal resource identifier. This internal resource identifier can be
- an eos fileid, because eos can look up files by id
- an inode if the filesystem and the storage driver support looking up by inode
- a path if the storage driver has no way of looking up files by id.
  - In this case other mechanisms like inotify, kernel audit or a fuse overlay might be used to keep the paths up to date.
  - to prevent excessive writes when deep folders are renamed a reverse map might be used: it will map the `uuid` to `<parentuuid>:<childname>`, in order to trade writes for reads
  - as a fallback a sync job can read the file id from the metadata of the resources and populate the uuid to internal id map.

The TUS upload can take metadata, for PUT we might need a header.