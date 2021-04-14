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

The below diagram shows the ocis services and the contained reva services within as dashed boxes. In general:
1. A request comes in at the proxy and is authenticated using oidc.
2. It is forwarded to the ocis frontend which handles ocs and ocdav requests by talking to the reva gateway using the CS3 API.
3. The gateway acts as a facade to the actual CS3 services: storage providers, user providers, group providers and sharing providers.

{{< svg src="extensions/storage/static/overview.drawio.svg" >}}

The dashed lines in the diagram indicate requests that are made to authenticate requests or lookup the storage provider:
1. After authenticating a request the proxy may either use the CS3 `userprovider` or the accounts service to fetch the user information that will be minted into the `x-access-token`.
2. The gateway will verify the JWT signature of the `x-access-token` or try to authenticate the request itself, eg. using a public link token.

{{< hint warning >}}
The bottom part is lighter because we will deprecate it in favor of using only the CS3 user and group providers after moving some account functionality into reva and glauth. The metadata storage is not registered in the reva gateway to prevent seperate user metadata necessary for running the service from data that is being served directly.
{{< /hint >}}

## Terminology

Communication is hard. And clear communication is even harder. You may encounter the following terms throughout the documentation, in the code or when talking to other developers. Just keep in mind that whenever you hear or read *storage*, that term needs to be clarified, because on its own it is too vague. PR welcome.

### Resources
A *resource* is a logical concept. It can be have different [types](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceType)
- an actual *file*
- a *container*, eg. a folder or bucket
- a *symlink*, or
- a *reference* which can point to a resource in another *storage provider*

### References

A [*reference*](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Reference) is a logical concept. It identifies a *resource* and consists of either
- a *path* based reference, used to identify a *resource* in the *namespace* of a *storage provider*. It must start with a `/`.
- an [*id* based reference](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId), uniquely identifying a *resounce* in the *namespace* of a *storage provider*. It consists of a `storage provider id` and an `opaque id`. The `storage provider id` must NOT start with a `/`.

{{< hint info >}}
The `/` is important because currenty the static *storage registry* uses a map to look up which *storage provider* is responsible for the resource. Paths must be prefixed with `/` so there can be no collisions between paths and storage provider ids in the same map.
{{< /hint >}}

{{< hint warning >}}
**Proposed Change**

A *reference* is a logical concept. It identifies a *resource* and consists of
a *root* and a *path* relative to that *root*. 
A *root* is a `storage space id` and a `logical uuid` of a resource inside the *storage space*:
```
<storages_pace>:<logical_uuid>:<relative_path>
`------------root------------´ `----path-----´
```
Both, *root* and *path*, are optional. The `storage space id` should be a uuid to prevent references from breaking when a *user* or *storage space* gets renamed. These are all valid references:

| name | description |
|------|-------------|
|`::/users/alice/projects/foo` | follow this path in the  global namespace (path based reference)|
|`home-alice::projects/foo` | in the *storage space* `home-alice`, start at the *root*, follow `projects/foo` |
|`home-alice:c3cf23bb-8f47-4719-a150-1d25a1f6fb56:foo` | in the *storage space* `home-alice`, start at the resource identified by the `logical uuid` `c3cf23bb-8f47-4719-a150-1d25a1f6fb56`, follow `projects/foo` |
|`home-alice:56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:` | in the *storage space* `home-alice`, start at the resource identified by the `logical uuid` `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a`, no relativ path to follow |
|`ee1687e5-ac7f-426d-a6c0-03fed91d5f62:56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:` | in the *storage space* `ee1687e5-ac7f-426d-a6c0-03fed91d5f62`, start at the resource identified by the `logical uuid` `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a`, no relativ path to follow (id based reference) |

They all reference the same *resource* (if the `logical uuid` of `projects` is 
`c3cf23bb-8f47-4719-a150-1d25a1f6fb56` and the `logical uuid` of `foo` is `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a`).
To move *storage spaces* between *storage providers*, the `logical uuid` of a file must be settable by the system. This also allows importing file metadata during a migration and restoring correct file ids when restoring a backup.
{{< /hint >}}


### Namespaces
A *namespace* is a set of paths with a common prefix. There is a global namespace when making requests against the CS3 *gateway* and a local namespace for every *storage provider*. Currently, oCIS mounts these *storage providers* out of the box:

| mountpoint | served namespace |
|------------|------------------|
| `/home`    | currently logged in users home |
| `/users`   | all users, used to access collaborative shares |
| `/public`  | publicly shared files, used to access public links |


{{< hint warning >}}
We plan to serve all shares the current user has access to under `/shares`.
{{< /hint >}}

### Storage Drivers

A *storage driver* implements access to a *storage system*:

It maps the *path* and *id* based CS3 *references* to an appropriate *storage system* specific reference, eg:
- eos file ids
- posix inodes or paths
- deconstructed filesystem nodes

### Storage Provider

A *storage provider* manages *resources* identified by a *reference*
by accessing a *storage system* with a *storage driver*.

{{< svg src="extensions/storage/static/storageprovider.drawio.svg" >}}

{{< hint warning >}}
**Proposed Change**
A *storage provider* manages multiple *storage spaces*
by accessing a *storage system* with a *storage driver*.

{{< /hint >}}
{{< svg src="extensions/storage/static/storageprovider-spaces.drawio.svg" >}}
{{< hint warning >}}
By making *storage providers* aware of *storage spaces* we can get rid of the current `enablehome` flag / hack in reva. Furthermore, provisioning a new *storage space* becomes a generic operation, regardless if a new user home or a new project space needs to be provisioned. 
{{< /hint >}}

### Storage Registries

A *storage registry* manages the global *namespace*:
it is used by the *gateway*
to look up `address` and `port` of the *storage provider*
that should handle a *reference*.

{{< svg src="extensions/storage/static/storageregistry.drawio.svg" >}}

{{< hint warning >}}
**Proposed Change**
A *storage registry* manages the *namespace* for a *user*:
it is used by the *gateway*
to look up `address` and `port` of the *storage provider*
that is currently serving a *storage space*.
{{< /hint >}}
{{< svg src="extensions/storage/static/storageregistry-spaces.drawio.svg" >}}
{{< hint warning >}}
By making *storage registries* aware of *storage spaces* we can query them for a listing of all *storage spaces* a user has access to. Including his home, received shares, project folders or group drives. See [Add draft of adr for spaces API ocis#1827](https://github.com/owncloud/ocis/pull/1827) for more info.
{{< /hint >}}

### Storage Spaces
A *storage space* is a logical concept:
it is a tree of *resources*
with a single *owner* (*user* or *group*), 
a *quota* and *permissions*, identified by a `storage space id`.

{{< svg src="extensions/storage/static/storagespace.drawio.svg" >}}

Examples would be every users home storage space, project storage spaces or group storage spaces. While they all serve different purposes and may or may not have workflows like anti virus scanning enabled, we need a way to identify and manage these subtrees in a generic way. By creating a dedicated concept for them this becomes easier and literally makes the codebase cleaner.

Finally, a logical `storage space id` is not tied to a specific storage provider. When the *storage driver* supports it we can import existing files including their fileid, which makes it possible to move storage spaces between storage spaces to implement storage classes eg. with or without archival, workflows, on SSDs or HDDs.