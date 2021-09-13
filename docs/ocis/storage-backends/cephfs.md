---
title: "cephfs"
date: 2021-09-13T15:36:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/storage-backends/
geekdocFilePath: cephfs.md
---

{{< toc >}}

oCIS intends to make the aspects of existing storage systems available as transparently as possible, but the static sync algorithm of the desktop client relies on some form of recursive change time propagation on the server side to detect changes. While this can be bolted on top of existing file systems with inotify, the kernel audit or a fuse based overlay filesystem, a storage system that already implements this aspect is preferable. Aside from EOS, cephfs supports a recursive change time that oCIS can use to calculate an etag for the webdav API.

## Development

The cephfs development happens in a [reva branch](https://github.com/cs3org/reva/pull/1209) and is currently driven by CERN. 

## Architecture

In the original approach the driver was based on the localfs driver, relying on a locally mounted cephfs. It would interface with it using the POSIX apis. This has been changed to direct Ceph API access using https://github.com/ceph/go-ceph. It allows using the ceph admin APIs to create subvolumes for user homes and maintain a file id to path mapping using symlinks.

It also uses the `.snap` folder built into Ceph to provide versions.

Trash is not implemented, as cephfs has no native recycle bin.

## Future work
- The spaces concept matches subvolumes, implement the CreateStorageSpace call with that, keep track of the list of storage spaces using symlings, like for the id based lookup
- The Share manager needs a persistence layer.
  - currently we persist using a json file. An sqlite db would be more robust.
  - As it basically provides two lists, *shared with me* and *shared with others* we could persist this directly on cephfs!
  - To allow deprovisioning a user the data should by sharded by userid.
  - Backups are then done using snapshots.