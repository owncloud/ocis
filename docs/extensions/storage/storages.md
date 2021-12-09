---
title: "Storages"
date: 2020-04-27T18:46:00+01:00
weight: 37
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: storages.md
---

## Storage commands

`storage` has multiple storage provider commands to preconfigure different default configurations for the reva *storage provider* service. While you could rerun `storage storage-oc` multiple times with different flags to get multiple instances we are giving the different commands the necessary default configuration to allow the `ocis` binary to simply start them and not deal with configuration.

## Storage providers

To manage the file tree oCIS uses *storage providers* that are accessing the underlying storage using a *storage driver*. The driver can be used to change the implementation of a storage aspect to better reflect the actual underlying storage capabilities. As an example a move operation on a POSIX filesystem ([theoretically](https://danluu.com/deconstruct-files/)) is an atomic operation. When trying to implement a file tree on top of S3 there is no native move operation that can be used. A naive implementation might fall back on a COPY and DELETE. Some S3 implementations provide a COPY operation that uses an existing key as the source, so the file at least does not need to be reuploaded. In the worst case scenario, which is renaming a folder with hundreds of thousands of objects, a reupload for every file has to be made. Instead of hiding this complexity a better choice might be to disable renaming of files or at least folders on S3. There are however implementations of filesystems on top of S3 that store the tree metadata in dedicated objects or use a completely different persistence mechanism like a distributed key value store to implement the file tree aspect of a storage.


{{< hint info >}}
While the *storage provider* is responsible for managing the tree, file up- and downloads are delegated to a dedicated *data provider*. See below.
{{< /hint >}}

## Storage aspects
A lot of different storage technologies exist, ranging from general purpose file systems with POSIX semantics to software defined storage with multiple APIs. Choosing any of them is making a tradeoff decision. Or, if a storage technology is already in place it automatically predetermines the capabilities that can be made available. *Not all storage systems are created equal.*

Unfortunately, no POSIX filesystem natively supports all storage aspects that ownCloud 10 requires:


### A hierarchical file tree
An important aspect of a filesystem is organizing files and directories in a file hierarchy, or tree. It allows you to create, move and delete nodes. Beside the name a node also has well known metadata like size and mtime that are persisted in the tree as well.

{{< hint info >}}
**Folders are not directories**
There is a difference between *folder* and *directory*: a *directory* is a file system concept. A *folder* is a metaphor for the concept of a physical file folder. There are also *virtual folders* or *smart folders* like the recent files folder which are no file system *directories*. So, every *directory* and every *virtual folder* is a *folder*, but not every *folder* is a *directory*. See [the folder metaphor in wikipedia](https://en.wikipedia.org/wiki/Directory_(computing)#Folder_metaphor). Also see the activity history below.
{{< /hint >}}

#### Id based lookup
While traditionally nodes in the tree are reached by traversing the path the tree persistence should be prepared to look up a node by an id. Think of an inode in a POSIX filesystem. If this operation needs to be cached for performance reasons keep in mind that cache invalidation is hard and crawling all files to update the inode to path mapping takes O(n), not O(1).

#### ETag propagation
For the state based sync a client can discover changes by recursively descending the tree and comparing the ETag for every node. If the storage technology supports propagating ETag changes up the tree, only the root node of a tree needs to be checked to determine if a discovery needs to be started and which nodes need to be traversed. This allows using the storage technology itself to persist all metadata that is necessary for sync, without additional services or caches.

#### Subtree size accounting
The tree can keep track of how many bytes are stored in a folder. Similar to ETag propagation a change in file size is propagated up the hierarchy.

{{< hint info >}}
**ETag and Size propagation**
When propagating the ETag (mtime) and size changes up the tree the question is where to stop. If all changes need to be propagated to the root of a storage then the root or busy folders will become a hotspot. There are two things to keep in mind: 1. propagation only happens up to the root of a single space (a user private drive or a single group drive), 2. no cross storage propagation. The latter was used in oc10 to let clients detect when a file in a received shared folder changed. This functionality is moving to the storage registry which caches the ETag for every root so clients can discover if and which storage changed.
{{< /hint >}}

#### Rename
Depending on the underlying storage technology some operations may either be slow, up to a point where it makes more sense to disable them entirely. One example is a folder rename: on S3 a *simple* folder rename translates to a copy and delete operation for every child of the renamed folder. There is an exception though: this restriction only applies if the S3 storage is treated like a filesystem, where the keys are the path and the value is the file content. There are smarter ways to implement file systems on top of S3, but again: there is always a tradeoff.

{{< hint info >}}
**S3 has no rename**
Technically, [S3 has no rename operation at all](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/examples-s3-objects.html#copy-object). By design, the location of the value is determined by the key, so it always has to do a copy and delete. Another example is the [redis RENAME operation](https://redis.io/commands/rename): while being specified as O(1) it *executes an implicit DEL operation, so if the deleted key contains a very big value it may cause high latency...*
{{< /hint >}}

#### Arbitrary metadata persistence
In addition to well known metadata like name size and mtime, users might be able to add arbitrary metadata like tags, comments or [dublin core](https://en.wikipedia.org/wiki/Dublin_Core). In POSIX filesystems this maps to extended attributes.

### Grant persistence
The CS3 API uses grants to describe access permissions. Storage systems have a wide range of permissions granularity and not all grants may be supported by every storage driver. POSIX ACLs for example have no expiry. If the storage system does not support certain grant properties, e.g. expiry, then the storage driver may choose to implement them in a different way. Expiries could be persisted in a different way and checked periodically to remove the grants. Again: every decision is a tradeoff.

### Trash persistence
After deleting a node the storage allows listing the deleted nodes and has an undo mechanism for them.

### Versions persistence
A user can restore a previous version of a file.

{{< hint info >}}
**Snapshots are not versions**
Modern POSIX filesystems support snapshotting of volumes. This is different from keeping track of versions to a file or folder, but might be another implementation strategy for a storage driver to allow users to restore content.
{{< /hint >}}

### Activity History
The storage keeps an activity history, tracking the different actions that have been performed. This does not only include file changes but also metadata changes like renames and permission changes.

## Storage drivers

Reva currently has four storage driver implementations that can be used for *storage providers* an well as *data providers*.

### Local Storage Driver

The *minimal* storage driver for a POSIX based filesystem. It literally supports none of the storage aspect other than basic file tree management. Sharing can - to a degree - be implemented using POSIX ACLs.

- tree provided by a POSIX filesystem
  - inefficient path by id lookup, currently uses the file path as id, so ids are not stable
    - can store a uuid in extended attributes and use a cache to look them up, similar to the ownCloud driver
  - no native ETag propagation, five options are available:
    - built in propagation (changes bypassing ocis are not picked up until a rescan)
    - built in inotify (requires 48 bytes of RAM per file, needs to keep track of every file and folder)
    - external inotify (same RAM requirement, but could be triggered by external tools, e.g. a workflow engine)
    - kernel audit log (use the linux kernel audit to capture file events on the storage and offload them to a queue)
    - fuse filesystem overlay
  - no subtree accounting, same options as for ETag propagation
  - efficient rename
  - arbitrary metadata using extended attributes
- grant persistence
  - using POSIX ACLs
    - requires an LDAP server to make guest accounts available in the OS
      - oCIS has glauth which contains all users
      - an existing LDAP could be used if guests ar provisioned in another way
  - using extended attributes to implement expiry or sharing that does not require OS level integration
  - fuse filesystem overlay
- no native trash
  - could use the [The FreeDesktop.org Trash specification](https://specifications.freedesktop.org/trash-spec/trashspec-latest.html)
  - fuse filesystem overlay
- no native versions, multiple options possible
  - git for folders
  - rcs for single files
  - rsnapshot for hourly / daily / weekly / monthly backups ... but this is not versioning as known from oc10
  - design new freedesktop spec, basically what is done in oc10 without the limitations or borrow ideas from the freedesktop trash spec
  - fuse filesystem overlay

To provide the other storage aspects we plan to implement a FUSE overlay filesystem which will add the different aspects on top of local filesystems like ext4, btrfs or xfs. It should work on NFSv45 as well, although NFSv4 supports RichACLs and we will explore how to leverage them to implement sharing at a future date. The idea is to use the storages native capabilities to deliver the best user experience. But again: that means making the right tradeoffs.

### OwnCloud Storage Driver

This is the current default storage driver. While it implements the file tree (using redis, including id based lookup), ETag propagation, trash, versions and sharing (including expiry) using the data directory layout of ownCloud 10 it has [known limitations](https://github.com/owncloud/core/issues/28095) that cannot be fixed without changing the actual layout on disk.

To setup it up properly in a distributed fashion, the storage-home and the storage-oc need to share the same underlying FS. Their "data" counterparts also need access to the same shared FS.
For a simple docker-compose setup, you can create a volume which will be used by the "storage-storage-home", "storage-storage-home-data", "storage-storage-oc" and "storage-storage-oc-data" containers. Using the `owncloud/ocis` docker image, the volume would need to be hooked in the `/var/lib/ocis` folder inside the containers.

- tree provided by a POSIX filesystem
  - file layout is mapped to the old ownCloud 10 layout
    - the root of tree for a user on disk is prefixed with `/path/to/data/<username>/files/`
  - efficient path by id lookup
    - all files and folders get assigned a uuid in the extended attributes
    - when starting the storage provider it will walk all files to populate a redis kv store for uuid to path lookup
    - slow to boot trees with lots of nodes
  - build in ETag propagation
    - ETags are calculated based on mtime
    - mtime is propagated by the storage driver
    - changes bypassing ocis are not picked up until a restart of the storage provider
  - no subtree accounting, same options as for local storage
  - efficient rename
    - TODO [update the kv store for path lookup](https://github.com/cs3org/reva/issues/985), this is an O(n) operation
  - arbitrary metadata using extended attributes
- grant persistence
  - using custom ACLs that are stored as extended attributes
    - a grant corresponds to one extended attribute of 40-100 bytes, effectively limiting the number of shares to ~100-40
    - extended attributes have varying limitations, based on the underlying filesystem
      - the linux kernel imposes a limit of 255bytes per name and 64KiB per value
      - ext2/3/4: total bytes for all attributes of a file is limited to 4KiB (a filesystem block)
      - xfs: limit of 64KiB per value
      - btrfs: total bytes used for the name, value, and implementation overhead bytes 16KiB (the default filesystem nodesize value)
  -  does not require OS level integration
- built in trash
  - trashed files are moved to `/path/to/data/<username>/files_trashbin/`
  - trashed files are appended a timestamp `.d<unixtime>`, which [breaks trashing of files that reach the filesystems specific name limit](https://github.com/owncloud/core/issues/28095)
- built in versions
  - file versions are stored in `/path/to/data/<username>/files_versions/`
  - file versions are appended a timestamp `.d<unixtime>`, which [breaks versioning of files that reach the filesystems specific name limit](https://github.com/owncloud/core/issues/28095)

### EOS Storage Driver

The CERN eos storage has evolved with ownCloud and natively supports id based lookup, ETag propagation, subtree size accounting, sharing, trash and versions. To use it you need to change the default configuration of the `storage storage-home` command (or have a look at the Makefile ̀ eos-start` target):

```
export STORAGE_HOME_DRIVER=eos
export STORAGE_DRIVER_EOS_NAMESPACE=/eos
export STORAGE_DRIVER_EOS_MASTER_URL="root://eos-mgm1.eoscluster.cern.ch:1094"
export STORAGE_DRIVER_EOS_ENABLE_HOME=true
export STORAGE_DRIVER_EOS_LAYOUT="dockertest/{{.Username}}"
```

Running it locally also requires the `eos` and `xrootd` binaries. Running it using `make eos-start` will use CentOS based containers that already have the necessary packages installed.

{{< hint info >}}
Pull requests to add explicit `storage storage-(s3|custom|...)` commands with working defaults are welcome.
{{< /hint >}}

### S3 Storage Driver

A naive driver that treats the keys in an S3 capable storage as `/` delimited path names. While it does not support MOVE or ETag propagation it can be used to read and write files. Better integration with native capabilities like versioning is possible but depends on the Use Case. Several storage solutions that provide an S3 interface also support some form of notifications that can be used to implement ETag propagation.

## Data Providers

Clients using the CS3 API use an [InitiateFileDownload](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.InitiateFileDownloadRequest) and ]InitiateUpload](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.InitiateFileUploadRequest) request at the [storage gateway](https://cs3org.github.io/cs3apis/#cs3.gateway.v1beta1.GatewayAPI) to obtain a URL endpoint that can be used to either GET the file content or upload content using the resumable [tus.io](https://tus.io) protocol.

The *data provider* uses the same *storage driver* as the *storage provider* but can be scaled independently.

The dataprovider allows uploading the file to a quarantine area where further data analysis may happen before making the file accessible again. One use case for this is anti virus scanning for files coming from untrusted sources.

## Future work

### FUSE overlay filesystem
We are planning to further separate the concerns and use a local storage provider with a FUSE filesystem overlaying the actual POSIX storage that can be used to capture deletes and writes that might happen outside of ocis/reva.

It would allow us to extend the local storage driver with missing storage aspects while keeping a tree like filesystem that end users are used to see when sshing into the machine.

### Upload to Quarantine area
Antivirus scanning of random files uploaded from untrusted sources and executing metadata extraction or thumbnail generation should happen in a sandboxed system to prevent malicious users from gaining any information about the system. By spawning a new container with access to only the uploaded data we can further limit the attack surface.
