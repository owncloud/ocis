---
title: "16. Storage for Files Metadata"
weight: 16
date: 2022-03-02T00:00:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0016-files-metadata.md
---

* Status: superseded by [ADR-0024]({{< ref "0024-msgpack-metadata.md" >}})
* Deciders: [@butonic](https://github.com/butonic), [@dragotin](https://github.com/dragotin), [@micbar](https://github.com/micbar), [@c0rby](https://github.com/c0rby)
* Date: 2022-02-04

## Context and Problem Statement

In addition to the file content we need to store metadata which is attached to a file. Metadata describes additional properties of a file. These properties need to be stored as close as possible to the file content to avoid inconsistencies. Metadata are key to workflows and search. We consider them as an additional value which enhances the file content.

## Decision Drivers

* Metadata will become more important in the future
* Metadata are key to automated data processing
* Metadata storage should be as close as possible to the file content
* Metadata should be always in sync with the file content

## Considered Options

* Database
* Extended file attributes
* Metadata file next to the file content
* Linked metadata in separate file

## Decision Outcome

Chosen option: "Extended File Attributes", because we guarantee the consistency of data and have arbitrary simple storage mechanism.

### Positive Consequences

* Metadata is always attached to the file itself
* We can store arbitrary key/values
* No external dependencies are needed

### Negative consequences

* The storage inside extended file attributes has limits
* Changes to extended attributes are not atomic and need file locks

## Pros and Cons of the Options <!-- optional -->

### Database or Key-Value Store

Use a Database or an external key/value store to persist metadata.

* Good, because it scales well
* Good, because databases provide efficient lookup mechanisms
* Bad, because the file content and the metadata could run out of sync
* Bad, because a storage backup doesn't cover the file metadata

### Extended File Attributes

Extended File Attributes allow storing arbitrary properties. There are 4 namespaces `user`, `system`, `trusted` and `security`. We can safely use the `user` namespace. An example attribute name would be `user.ocis.owner.id`. The linux kernel has length limits on attribute names and values.

From Wikipedia on [Extended file attributes](https://en.wikipedia.org/wiki/Extended_file_attributes#Linux):

> The Linux kernel allows extended attribute to have names of up to 255 bytes and values of up to 64 KiB,[14] as do XFS and ReiserFS, but ext2/3/4 and btrfs impose much smaller limits, requiring all the attributes (names and values) of one file to fit in one “filesystem block” (usually 4 KiB). Per POSIX.1e,[citation needed] the names are required to start with one of security, system, trusted, and user plus a period. This defines the four namespaces of extended attributes.

* Good, because metadata is stored in the filesystem
* Good, because consistency is easy to maintain
* Good, because the data is attached to the file and survives file operations like copy and move
* Good, because a storage backup also covers the file metadata
* Bad, because we could hit the filesystem limit
* Bad, because changes to extended attributes are not atomic

### Metadata File

We could store metadata in a metadata file next to the file content which has a structured content format like .json, .yaml or .toml. That would give us more space to store bigger amounts of metadata.

* Good, because there are no size limits
* Good, because there is more freedom to the content format
* Good, because a storage backup also covers the file metadata
* Bad, because it doubles the amount of read / write operations
* Bad, because it needs additional measures against concurrent overwriting changes

### Link metadata with an id in the extended attributes

To link metadata to file content a single extended attribute with a file id (unique per storage space) is sufficient. This would also allow putting metadata in better suited storage systems like SQLite or a key value store.

* Good, because it avoids extended attribute limits
* Good, because the same mechanism could be used to look up files by id, when the underlying filesystem is an existing POSIX filesystem.
* Bad, because backup needs to cover the metadata as well. Could be mitigated by sharing metadata per space and doing space wide snapshots.
* Bad, because it is a bit more effort to access it to read or index it.
