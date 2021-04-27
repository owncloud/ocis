---
title: "Terminology"
date: 2018-05-02T00:00:00+00:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: terminology.md
---

Communication is hard. And clear communication is even harder. You may encounter the following terms throughout the documentation, in the code or when talking to other developers. Just keep in mind that whenever you hear or read *storage*, that term needs to be clarified, because on its own it is too vague. PR welcome.

## Resources
A *resource* is a logical concept. Ressources can be of [different types](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceType):
- an actual *file*
- a *container*, e.g. a folder or bucket
- a *symlink*, or
- a *reference* which can point to a resource in another *storage provider*

## References

A [*reference*](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Reference) is a logical concept. It identifies a *resource* and consists of either
- a *path* based reference, used to identify a *resource* in the *namespace* of a *storage provider*. It must start with a `/`.
- an [*id* based reference](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId), uniquely identifying a *resource* in the *namespace* of a *storage provider*. It consists of a `storage provider id` and an `opaque id`. The `storage provider id` must NOT start with a `/`.

{{< hint info >}}
The `/` is important because currently the static *storage registry* uses a map to look up which *storage provider* is responsible for the resource. Paths must be prefixed with `/` so there can be no collisions between paths and storage provider ids in the same map.
{{< /hint >}}


{{< hint warning >}}
### Alternative: reference triple ####
A *reference* is a logical concept. It identifies a *resource* and consists of
a `storage_space`, a `<root_id>` and a `<path>`
```
<storage_space>!<root_id>:<path>
```
While all components are optional, only three cases are used:
| format | example | description |
|-|-|-|
| `!:<absolute_path>` | `!:/absolute/path/to/file.ext` | absolute path | 
| `<storage_space>!:<relative_path>` | `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!:path/to/file.ext` | path relative to the root of the storage space | 
| `<storage_space>!<root>:<relative_path>` | `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!c3cf23bb-8f47-4719-a150-1d25a1f6fb56:to/file.ext` | path relative to the specified node in the storage space, used to reference resources without disclosing parent paths |

`<storage_space>` should be a uuid to prevent references from breaking when a *user* or *storage space* gets renamed. But it can also be derived from a migration of an oc10 instance by concatenating an instance identifier and the numeric storage id from oc10, e.g. `oc10-instance-a$1234`.

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


The `:`, `!` and `$` are chosen deliberately from the set of [RFC3986 sub delimiters](https://tools.ietf.org/html/rfc3986#section-2.2), so they can be used in URLs without having to being encoded. In some cases, a delimiter can be left out if a component is not set:
| reference | interpretation |
|-|-|
| `/absolute/path/to/file.ext` | absolute path, all delimiters omitted |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!path/to/file.ext` | relative path in the given storage space, root delimiter `:` omitted |
| `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:to/file.ext` | relative path in the given root node, storage space delimiter `!` omitted |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62!56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:` | node id in the given storage space, `:` must be present |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62` | root of the storage space, all delimiters omitted, can be distinguished by the `/` |

{{< /hint >}}

## Storage Drivers

A *storage driver* implements access to a *storage system*:

It maps the *path* and *id* based CS3 *references* to an appropriate *storage system* specific reference, e.g.:
- eos file ids
- posix inodes or paths
- deconstructed filesystem nodes

## Storage Provider

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
By making *storage providers* aware of *storage spaces* we can get rid of the current `enablehome` flag / hack in reva. Furthermore, provisioning a new *storage space* becomes a generic operation, regardless of the need of provisioning a new user home or a new project space.
{{< /hint >}}

## Storage Registries

A *storage registry* manages the global *namespace*:
It is used by the *gateway*
to look up `address` and `port` of the *storage provider*
that should handle a *reference*.

{{< svg src="extensions/storage/static/storageregistry.drawio.svg" >}}

{{< hint warning >}}
**Proposed Change**
A *storage registry* manages the *namespace* for a *user*:
It is used by the *gateway*
to look up `address` and `port` of the *storage provider*
that is currently serving a *storage space*.
{{< /hint >}}
{{< svg src="extensions/storage/static/storageregistry-spaces.drawio.svg" >}}
{{< hint warning >}}
By making *storage registries* aware of *storage spaces* we can query them for a listing of all *storage spaces* a user has access to. Including his home, received shares, project folders or group drives. See [a WIP PR for spaces in the oCIS repo (#1827)](https://github.com/owncloud/ocis/pull/1827) for more info.
{{< /hint >}}

## Storage Spaces
A *storage space* is a logical concept:
It is a tree of *resources*
with a single *owner* (*user* or *group*), 
a *quota* and *permissions*, identified by a `storage space id`.

{{< svg src="extensions/storage/static/storagespace.drawio.svg" >}}

Examples would be every user's home storage space, project storage spaces or group storage spaces. While they all serve different purposes and may or may not have workflows like anti virus scanning enabled, we need a way to identify and manage these subtrees in a generic way. By creating a dedicated concept for them this becomes easier and literally makes the codebase cleaner. A *storage registry* then allows listing the properties of storage spaces, e.g. free space, quota, owner, syncable, root etag, uploed workflow steps, ...

Finally, a logical `storage space id` is not tied to a specific storage provider. If the *storage driver* supports it, we can import existing files including their *file id*, which makes it possible to move storage spaces between storage spaces to implement storage classes, e.g. with or without archival, workflows, on SSDs or HDDs.

## Shares
*To be clarified: we are aware that storage spaces may be to 'heavywheight' for ad hoc sharing with groups. That being said there is no technical reason why group shares should notd be treated like storage spaces that users can provision themselves. They would share the quota with the users home storage space and the share initiator would be the sole owner but the mechanism of treating a share like a new storage space would be the same. This obviously also extends to user shares and even file indvidual shares that would be wrapped in a virtual collection. Something new that would become possible would be collections of arbitrary files in a single storage space, e.g. the ten best pictures from a large album.*