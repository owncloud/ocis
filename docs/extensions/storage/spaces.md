---
title: "Spaces"
date: 2018-05-02T00:00:00+00:00
weight: 3
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: spaces.md
---

{{< hint warning >}}

The current implementation in oCIS might not yet fully reflect this concept. Feel free to add links to ADRs, PRs and Issues in short warning boxes like this.

{{< /hint >}}

## Storage Spaces
A storage *space* is a logical concept. It organizes a set of [*resources*]({{< ref "#resources" >}}) in a hierarchical tree. It has a single *owner* (*user* or *group*), 
a *quota*, *permissions* and is identified by a `storage space id`.

{{< svg src="extensions/storage/static/storagespace.drawio.svg" >}}

Examples would be every user's personal storage *space*, project storage *spaces* or group storage *spaces*. While they all serve different purposes and may or may not have workflows like anti virus scanning enabled, we need a way to identify and manage these subtrees in a generic way. By creating a dedicated concept for them this becomes easier and literally makes the codebase cleaner. A storage [*Spaces Registry*]({{< ref "./spacesregistry.md" >}}) then allows listing the capabilities of storage *spaces*, e.g. free space, quota, owner, syncable, root etag, upload workflow steps, ...

Finally, a logical `storage space id` is not tied to a specific [*spaces provider*]({{< ref "./spacesprovider.md" >}}). If the [*storage driver*]({{< ref "./storagedrivers.md" >}}) supports it, we can import existing files including their `file id`, which makes it possible to move storage *spaces* between [*spaces providers*]({{< ref "./spacesprovider.md" >}}) to implement storage classes, e.g. with or without archival, workflows, on SSDs or HDDs.

## Shares
*To be clarified: we are aware that [*storage spaces*]({{< ref "#storage-spaces" >}}) may be too 'heavywheight' for ad hoc sharing with groups. That being said, there is no technical reason why group shares should not be treated like storage [*spaces*]({{< ref "#storage-spaces" >}}) that users can provision themselves. They would share the quota with the users home or personal storage [*space*]({{< ref "#storage-spaces" >}}) and the share initiator would be the sole owner. Technically, the mechanism of treating a share like a new storage [*space*]({{< ref "#storage-spaces" >}}) would be the same. This obviously also extends to user shares and even file individual shares that would be wrapped in a virtual collection. It would also become possible to share collections of arbitrary files in a single storage space, e.g. the ten best pictures from a large album.*

