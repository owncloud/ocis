---
title: "Proposed Changes"
date: 2018-05-02T00:00:00+00:00
weight: 18
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/storage
geekdocFilePath: proposedchanges.md
---

Some architectural changes still need to be clarified or changed. Maybe an ADR is in order for all of the below.

## Reva Gateway changes

## A dedicated shares storage provider

Currently, when a user accepts a share, a cs3 reference is created in the users `/home/shares` folder. This reference represents the mount point of a share and can be renamed, similar to the share jail in ownCloud 10. This spreads the metadata of a share in two places:
- the share is persisted in the *share manager*
- the mount point of a share is persisted in the home *storage provider*

Furthermore, the *gateway* treats `/home/shares` different than any other path: it will stat all children and calculate an etag to allow clients to discover changes in accepted shares. This requires the storage provider to cooperate and provide this special `/shares` folder in the root of a users home when it is accessed as a home storage. That is the origin of the `enable_home` config flag that needs to be implemented for every storage driver.

In order to have a single source of truth we need to make the *share manager* aware of the mount point. We can then move all the logic that aggregates the etag in the share folder to a dedicated *shares storage provider* that is using the *share manager* for persistence. The *shares storage provider* would provide a `/shares` namespace outside of `/home` that lists all accepted shares for the current user. As a result the storage drivers no longer need to have a `enable_home` flag that jails users into their home. The `/home/shares` folder would move outside of the `/home`. In fact `/home` will no longer be needed, because the home folder concept can be implemented as a space: `CreateHome` would create a `personal` space on the.

Work on this is done in https://github.com/cs3org/reva/pull/2023

{{< hint warning >}}
What about copy pasting links from the browser? Well this storage is only really needed to have a path to ocm shares that actually reside on other instances. In the UI the shares would be listed by querying a *share manager*. It returns ResourceIds, which can be stated to fetch a path that is then accessible in the CS3 global namespace. Two caveats:
- This only works for resources that are actually hosted by the current instance. For those it would leak the parent path segments to a shared resource.
- For accepted OCM shares there must be a path in the [*CS3 global namespace*]({{< ref "./namespaces.md#cs3-global-namespaces" >}}) that has to be the same for all users, otherwise they cannot copy and share those URLs.
{{< /hint >}}

### The gateway should be responsible for path transformations

Currently, storage providers are aware af their mount point, coupling them tightly with the gateway.

Tracked in https://github.com/cs3org/reva/issues/578

Work is done in https://github.com/cs3org/reva/pull/1866

## URL escaped string representation of a CS3 reference

For the spaces concept we introduced the `/dav/spaces/` endpoint. It encodes a cs3 *reference* in a URL compatible way.
1. We can separate the path using a `/`: `/dav/spaces/<spaceid>/<path>`
2. The `spaceid` currently is a cs3 resourceid, consisting of `<storageid>` and `<opaqueid>`. Since the opaqueid might contain `/` e.g. for the local driver we have to urlencode the spaceid.

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

### Space id vs resource id vs storage id

We have `/dav/meta/<fileid>` where the `fileid` is a string that was returned by a PROPFIND or by the `/graph/v1.0/me/drives/` endpoint? That returns a space id and the root drive item which has an `id`

Does that `id` have a specific format? We currently concatenate as `<storageid>!<nodeid>`.

A request against `/dav/meta/fileid` will use the reva storage registry to look up a path.

What if the storage space is moved to another storage provider. This happens during a migration:

1. the current oc10 fileids need to be prefixed with at least the numeric storage id to shard them.

`123` becomes `instanceprefix$345!123` if we use a custom prefix that identifies an instance (so we can merge multiple instances into one ocis instance) and append  the numeric storageid `345`. The pattern is `<instanceprefix>$<numericstorageid>!<fileid>`.

Every `<instanceprefix>$<numericstorageid>` identifies a space.

- [ ] the owncloudsql driver can return these spaceids when listing spaces.

Why does it not work if we just use the fileid of the root node in the db?

Say we have a space with three resources:
`<instanceprefix>$<numericstorageid>!<fileid>`
`instanceprefix$345!1`
`instanceprefix$345!2`
`instanceprefix$345!3`

All users have moved to ocis and the registry contains a regex to route all `instanceprefix.*` references to the storageprovider with the owncloudsql driver. It is up to the driver to locate the correct resource by using the filecache table. In this case the numeric storage id is unnecessary.

Now we migrate the space `345` to another storage driver:
- the storage registry contains a new entry for `instanceprefix$345` to send all resource ids for that space to the new storage provider
- the new storage driver has to take into account the full storageid because the nodeid may only be unique per storage space.

If we now have to fetch the path on the `/dav/meta/` endpoint:
`/dav/meta/instanceprefix$345!1`
`/dav/meta/instanceprefix$345!2`
`/dav/meta/instanceprefix$345!3`

This would work because the registry always sees `instanceprefix$345` as the storageid.

Now if we use the fileids directly and leave out the numeric storageid:
`<instanceprefix>!<fileid>`
`instanceprefix!1`
`instanceprefix!2`
`instanceprefix!3`

This is the current `<storageid>!<nodeid>` format.

The reva storage registry contains a `instanceid` entry pointing to the storage provider with the owncloudsql driver.

Resources can be looked up because the oc_filecache has a unique fileid over all storages.

Now we again migrate the space `345` to another storage driver:
- the storage registry contains a new entry for `instanceprefix!1` so the storage space root now points to the new storage provider
- The registry needs to be aware of node ids to route properly. This is a no-go. We don't want to keep a cache of *all* nodeids in the registry. Only the root nodes of spaces.
- The new storage driver only has a nodeid which might collide with other nodeids from other storage spaces, e.g. when two instances are imported into one ocis instance. Although it would be possible to just set up two storage providers extra care would have to be taken to prevent nodeid collisions when importing a space.

If we now have to fetch the path on the `/dav/meta/` endpoint:
`/dav/meta/instanceprefix!1` would work because it is the root of a space
`/dav/meta/instanceprefix!2` would cause the gateway to poll all storage providers because the registry has no way to determine the responsible storage provider
`/dav/meta/instanceprefix!3` same

The problem is that without a part in the storageid that allows differentiating storage spaces we cannot route them individually.

Now, we could use the nodeid of the root of a storage space as the spaceid ... if it is a uuid. If it is numeric it needs a prefix to distinguish it from other spaces.
`<space-root-uuid>!<fileid>` would be easy for the decomposedfs.
eos might use numeric ids: `<eosprefix>$<space-root-fileid>!<fileid>`, but it needs a custom prefix to distinguish multiple eos instances.

Furthermore, when migrating spaces between storage providers we want to stay collision free, which is why we should recommend uuids.

All this has implications for the decomposedfs, because it needs to split the nodes per space to prevent them from colliding.
