---
title: "Terminology"
date: 2018-05-02T00:00:00+00:00
weight: 17
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: terminology.md
---

Communication is hard. And clear communication is even harder. You may encounter the following terms throughout the documentation, in the code or when talking to other developers. Just keep in mind that whenever you hear or read *storage*, that term needs to be clarified, because on its own it is too vague. PR welcome.

## Logical concepts

### Resources
A *resource* is the basic building block that oCIS manages. It can be of [different types](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceType):
- an actual *file*
- a *container*, e.g. a folder or bucket
- a *symlink*, or
- a [*reference*]({{< ref "#references" >}}) which can point to a resource in another [*storage provider*]({{< ref "#storage-providers" >}})

### References

A *reference* identifies a [*resource*]({{< ref "#resources" >}}). A [*CS3 reference*](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Reference) can carry a *path* and a [CS3 *resource id*](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId). The references come in two flavors: absolute and combined.
Absolute references have either the *path* or the *resource id* set:
- An absolute *path* MUST start with a `/`. The *resource id* MUST be empty.
- An absolute *resource id* uniquely identifies a [*resource*]({{< ref "#resources" >}}) and is used as a stable identifier for sharing. The *path* MUST be empty.
Combined references have both, *path* and *resource id* set:
- the *resource id* identifies the root [*resource*]({{< ref "#resources" >}})
- the *path* is relative to that root. It MUST start with `.`


### Storage Spaces
A *storage space* organizes a set of [*resources*]({{< ref "#resources" >}}) in a hierarchical tree. It has a single *owner* (*user* or *group*), 
a *quota*, *permissions* and is identified by a `storage space id`.

{{< svg src="extensions/storage/static/storagespace.drawio.svg" >}}

Examples would be every user's personal storage space, project storage spaces or group storage spaces. While they all serve different purposes and may or may not have workflows like anti virus scanning enabled, we need a way to identify and manage these subtrees in a generic way. By creating a dedicated concept for them this becomes easier and literally makes the codebase cleaner. A [*storage space registry*]({{< ref "#storage-space-registries" >}}) then allows listing the capabilities of [*storage spaces*]({{< ref "#storage-spaces" >}}), e.g. free space, quota, owner, syncable, root etag, upload workflow steps, ...

Finally, a logical `storage space id` is not tied to a specific [*storage provider*]({{< ref "#storage-providers" >}}). If the [*storage driver*]({{< ref "#storage-drivers" >}}) supports it, we can import existing files including their `file id`, which makes it possible to move [*storage spaces*]({{< ref "#storage-spaces" >}}) between [*storage providers*]({{< ref "#storage-providers" >}}) to implement storage classes, e.g. with or without archival, workflows, on SSDs or HDDs.

### Shares
*To be clarified: we are aware that [*storage spaces*]({{< ref "#storage-spaces" >}}) may be too 'heavywheight' for ad hoc sharing with groups. That being said, there is no technical reason why group shares should not be treated like [*storage spaces*]({{< ref "#storage-spaces" >}}) that users can provision themselves. They would share the quota with the users home [*storage space*]({{< ref "#storage-spaces" >}}) and the share initiator would be the sole owner. Technically, the mechanism of treating a share like a new [*storage space*]({{< ref "#storage-spaces" >}}) would be the same. This obviously also extends to user shares and even file indvidual shares that would be wrapped in a virtual collection. It would also become possible to share collections of arbitrary files in a single storage space, e.g. the ten best pictures from a large album.*


### Storage Space Registries

A *storage space registry* manages the [*namespace*]({{< ref "./namespaces.md" >}}) for a *user*: it is used by *clients* to look up storage spaces a user has access to, the `/dav/spaces` endpoint to access it via WabDAV, and where the client should mount it in the users personal namespace.

{{< svg src="extensions/storage/static/spacesregistry.drawio.svg" >}}


## Technical concepts

### Storage Drivers

A *storage driver* implements access to a [*storage system*]({{< ref "#storage-systems" >}}):

It maps the *path* and *id* based CS3 *references* to an appropriate [*storage system*]({{< ref "#storage-systems" >}}) specific reference, e.g.:
- eos file ids
- posix inodes or paths
- deconstructed filesystem nodes

### Storage Providers

A *storage provider* manages [*resources*]({{< ref "#resources" >}}) identified by a [*reference*]({{< ref "#references" >}})
by accessing a [*storage system*]({{< ref "#storage-systems" >}}) with a [*storage driver*]({{< ref "#storage-drivers" >}}).

{{< svg src="extensions/storage/static/storageprovider.drawio.svg" >}}

### Storage Registry

A *storage registry* manages the [*CS3 global namespace*]({{< ref "./namespaces.md#cs3-global-namespaces" >}}):
It is used by the *gateway*
to look up `address` and `port` of the [*storage provider*]({{< ref "#storage-providers" >}})
that should handle a [*reference*]({{< ref "#references" >}}).

{{< svg src="extensions/storage/static/storageregistry.drawio.svg" >}}

### Storage Systems
Every *storage system* has different native capabilities like id and path based lookups, recursive change time propagation, permissions, trash, versions, archival and more.
A [*storage provider*]({{< ref "#storage-providers" >}}) makes the storage system available in the CS3 API by wrapping the capabilities as good as possible using a [*storage driver*]({{< ref "#storage-drivers" >}}).
There migt be multiple [*storage drivers*]({{< ref "#storage-drivers" >}}) for a *storage system*, implementing different tradeoffs to match varying requirements.

### Gateways
A *gateway* acts as a facade to the storage related services. It authenticates and forwards API calls that are publicly accessible.

