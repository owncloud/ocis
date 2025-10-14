---
title: "13. Locking"
weight: 13
date: 2021-08-17T12:56:53+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0013-locking.md
---

- Status: accepted
- Deciders: [@hodyroff](https://github.com/hodyroff), [@pmaier1](https://github.com/pmaier1), [@jojowein](https://github.com/jojowein), [@dragotin](https://github.com/dragotin), [@micbar](https://github.com/micbar), [@tbsbdr](https://github.com/tbsbdr), [@wkloucek](https://github.com/wkloucek)
- Date: 2021-11-03

## Context and Problem Statement

At the time of this writing no locking mechanisms exists in oCIS / REVA for both directories and files. The CS3org WOPI server implements a file based locking in order to lock files. This ADR discusses if this approach is ok for the general availability of oCIS or if changes are needed.

## Decision Drivers

- Is the current situation acceptable for the GA
- Is locking needed or can we have oCIS / REVA without locking

## Considered Options

1. File based locking
2. No locking
3. CS3 API locking

## Decision Outcome

For the GA we chose option 2. Therefore we need to remove or disable the file based locking functionality of the CS3org WOPI server. The decision was taken because the current file based locking does not work on file-only shares. The current locking also does not guarantee exclusive access to a file since other parts of oCIS like the WebDAV API or other REVA services don't respect the locks.

After the GA we need to implement option 3.

## Pros and Cons of the Options

### File based locking

The CS3org WOPI server creates a `.sys.wopilock.<filename>.` and `.~lock.<filename>#` file when opening a file in write mode

**File based locking is good**, because:

- it is already implemented in the current CS3org WOPI server

**File based locking is bad**, because:

- lock files should be checked by all parties manipulating files (e.g. the WebDAV api)
- lock files can be deleted by everyone
- you can not lock files in a file-only share (you need a folder share to create a lock file besides the original file)

If we have file based locks, we can also sync them with e.g. the Desktop Client.

**Syncing lock files is good**: because

- native office applications can notice lock files by the WOPI server and vice versa (LibreOffice also creates `.lock.<filename>#` files)

**Syncing lock files is bad**, because:

- if lockfile is not deleted, no one can edit the file
- creating lock files in a folder shared with 2000000 users creates a lot of noise and pressure on the server (etag propagation, therefore oC Desktop sync client has an ignore rule for `.~lock.*` files)

### No locking

We remove or disable the file based locking of the CS3org WOPI server.

**No locking is good**, because:

- you don't need to release locks
- overwriting a file just creates a new version of it

**No locking is bad**, because:

- merging changes from different versions is a pain, since there is no way to calculate differences for most of the files (e.g. docx or xlsx files)
- no locking breaks the WOPI specs, as the CS3 WOPI server won't be capable to honor the WOPI Lock related operations

### CS3 API locking

- Add CS3 API for resource (files, directories) locking, unlocking and checking locks
  - locking always with timeout
  - lock creation is a "create-if-not-exists" operation
  - locks need to have arbitrary metadata (e.g. the CS3 WOPI server is stateless by storing information on / in the locks)
- Implement WebDAV locking using the CS3 API
- Implement Locking in storage drivers
- Change CS3 WOPI server to use CS3 API locking mechanism
- Optional: manual lock / unlock in ownCloud Web (who is allowed to unlock locks of another user?)

**CS3 API locking is good**, because:

- you can lock files on the actual storage (if the storage supports that -> storage driver dependent)
- you can lock files in ownCloud 10 when using the ownCloudSQL storage driver in the migration deployment (but oC10 Collabora / OnlyOffice also need to implement locking, to fully leverage that)
- clients can get the lock information via the api without ignoring / hiding lock file changes
- clients can use the lock information to lock the file in their context (e.g. via some file explorer integration)

**CS3 API locking is bad**, because:

- it needs to be defined and implemented, currently not planned for the GA
