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

In the original approach the driver was based on the localfs driver, relying on a locally mounted cephfs. It would interface with it using the POSIX apis. This has been changed to directly call the Ceph API using https://github.com/ceph/go-ceph. It allows using the ceph admin APIs to create subvolumes for user homes and maintain a file id to path mapping using symlinks.

## Implemented Aspects
The recursive change time built ino cephfs is used to implement the etag propagation expected by the ownCloud clients. This allows oCIS to pick up changes that have been made by external tools, bypassing any oCIS APIs. 

Like other filesystems cephfs uses inodes and like most other filesystems inodes are reused. To get stable file identifiers the current cephfs driver assigns every node a file id and maintains a custom fileid to path mapping in a system directory:
```
/tmp/cephfs $ tree -a
.
├── reva
│   └── einstein
│       ├── Pictures
│       └── welcome.txt
└── .reva_hidden
    ├── .fileids
    │   ├── 50BC39D364A4703A20C58ED50E4EADC3_570078 -> /tmp/cephfs/reva/einstein
    │   ├── 571EFB3F0ACAE6762716889478E40156_570081 -> /tmp/cephfs/reva/einstein/Pictures
    │   └── C7A1397524D0419B38D04D539EA531F8_588108 -> /tmp/cephfs/reva/einstein/welcome.txt
    └── .uploads
```

Versions are not file but snapshot based, a [native feature of cephfs](https://docs.ceph.com/en/latest/dev/cephfs-snapshots/). The driver maps entries in the native cephfs `.snap` folder to the CS3 api recycle bin concept and makes them available in the web UI using the versions sidebar. Snepshots cen be triggered by users themselves or on a schedule.

Trash is not implemented, as cephfs has no native recycle bin and instead relies on the snapshot functionality that can be triggered by end users. It should be possible to automatically create a snapshot before deleting a file. This needs to be explored.

Shares [are be mapped to ACLs](https://github.com/cs3org/reva/pull/1209/files#diff-5e532e61f99bffb5754263bc6ce75f84a30c6f507a58ba506b0b487a50eda1d9R168-R224) supported by cephfs. The share manager is used to persist the intent of a share and can be used to periodically verify or reset the ACLs on cephfs.

## Future work
- The spaces concept matches cephfs subvolumes. We can implement the CreateStorageSpace call with that, keep track of the list of storage spaces using symlinks, like for the id based lookup.
- The share manager needs a persistence layer.
- Currently we persist using a single json file.
- As it basically provides two lists, *shared with me* and *shared with others*, we could persist them directly on cephfs!
  - If needed for redundancy, the share manager can be run multiple times, backed by the same cephfs
  - To save disk io the data can be cached in memory, and invalidated using stat requests.
- A good tradeoff would be a folder for each user with a json file for each list. That way, we only have to open and read a single file when the user want's to list the shares.    
- To allow deprovisioning a user the data should by sharded by userid.
- For consistency over metadata any file blob data, backups can be done using snapshots.
- An example where einstein has shared a file with marie would look like this on disk:
```
/tmp/cephfs $ tree -a
.
├── reva
│   └── einstein
│       ├── Pictures
│       └── welcome.txt
├── .reva_hidden
│   ├── .fileids
│   │   ├── 50BC39D364A4703A20C58ED50E4EADC3_570078 -> /tmp/cephfs/reva/einstein
│   │   ├── 571EFB3F0ACAE6762716889478E40156_570081 -> /tmp/cephfs/reva/einstein/Pictures
│   │   └── C7A1397524D0419B38D04D539EA531F8_588108 -> /tmp/cephfs/reva/einstein/welcome.txt
│   └── .uploads
└── .reva_share_manager
    ├── einstein
    │   └── sharedWithOthers.json
    └── marie
        └── sharedWithMe.json
```
- The fileids should [not be based on the path](https://github.com/cs3org/reva/pull/1209/files#diff-eba5c8b77ccdd1ac570c54ed86dfa7643b6b30e5625af191f789727874850172R125-R127) and instead use a uuid that is also persisted in the extended attributes to allow rebuilding the index from scratch if necessary.