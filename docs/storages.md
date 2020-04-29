---
title: "Storages"
date: 2020-04-27T18:46:00+01:00
weight: 37
geekdocRepo: https://github.com/owncloud/ocis-reva
geekdocEditPath: edit/master/docs
geekdocFilePath: storages.md
---

## Storage commands

`ocis-reva` has multiple storage provider commands to preconfigure different default configurations for the reva *storage provider* service. While you could rerun `ocis-reva storage-oc` multiple times with different flags to get multiple instances we are giving the different commands the necessary default configuration to allow the `ocis` binary to simply start them and not deal with configuration.

## Storage providers

To manage the file tree ocis uses reva *storage providers* that are accessing the underlying storage using a *storage driver*. The driver can be used to change the implementation of a storage aspect to better reflect the actual underying storage capabilities. As an example a move operation on a POSIX filesystem ([theoretically](https://danluu.com/deconstruct-files/)) is an atomic operation. When trying to implement a file tree on top S3 there is no native move operation that can be used. A naive implementation might fall bak on a COPY and DELETE. Some S3 implementations provide a COPY operation that uses an existing key as the source, so the file at least does not need to be reuploaded. In the worst case scenario, the rename of a folder with hundreds of thousands of objects, a reupload for every file has to be made. Instead of hiding this complexity a better choice might be to disable renaming of files or at least folders on S3. There are however implemetations of filesystems on top of S3 that store the tree metadata in dedicated objects or use a completely different persistance mechanism like a distributed key value store to implement the file tree aspect of a storage.


{{< hint info >}}
While the *storage provider* is responsible for managing the tree, file up and download is delegated to a dedicated *data provider*. See below.
{{< /hint >}}

## Storage aspects

Unfortunately, no POSIX filesystem natively supports all storage aspects that ownCloud 10 requires:
- a hierarchical file tree
  - id based lookup
  - etag propagation
  - subtree size accounting (size of all files in a folder and all its sub folders)
- sharing
  - share expiry
- trash or undelete (trash can be done by wrapping rm)
- versions (only snapshots, which is a different concept)

A more extensive description of the storage aspects can be found in the [upstream documentation](https://reva.link/docs/concepts/storages/#aspects-of-storage-drivers)

## Storage drivers

Reva currently has four storage driver implementations that can be used for *storage providers* an well as *data providers*.

### Local Storage Driver

The *minimal* storage driver for a POSIX based filesystem. It literally supports none of the storage aspect other than basic file tree management. Sharing can - to a degree - be implemented using POSIX ACLs, which is scheduled after finishing the eos storage driver.

To provide the other storage aspects we plan to implement a FUSE overlay filesystem which will add the different aspects on top of local filesystems like ext4, btrfs or xfs. It should work on NFSv45 as well, although NFSv4 supports RichACLs and we will explore how to leverage them to implement sharing at a future date.

### OwnCloud Storage Driver

This is the current default storage driver. While it implements the file tree (using redis, including id based lookup), etag propagation, trash, versions and sharing (including expiry) using the data directory layout of ownCloud 10 it has [known limitations](https://github.com/owncloud/core/issues/28095) that cannot be fixed without changing the actual layout on disk.

We plan to deprecate it in favor of the local storage driver in combination with a FUSE based overlay filesystem when the migration path has been fully tested.

### EOS Storage Driver

The cern eos storage has evolved with ownCloud and natively supports id based lookup, etag propagation, subtree size accounting, sharing, trash and versions. To use it you need to change the default configuration of the `ocis-reva storage-home` command (or have a look at the Makefile ̀ eos-start` target):

```
export REVA_STORAGE_HOME_DRIVER=eos
export REVA_STORAGE_EOS_NAMESPACE=/eos
export REVA_STORAGE_EOS_MASTER_URL="root://eos-mgm1.eoscluster.cern.ch:1094"
export REVA_STORAGE_EOS_ENABLE_HOME=true
export REVA_STORAGE_EOS_LAYOUT="dockertest/{{.Username}}"
```

Running it locally also requires the `eos` and `xrootd` binaries. Running it using `make eos-start` will use CentOS based containers that already have the necessary packages installed.

{{< hint info >}}
Pull requests to add explicit `reva storage-(s3|custom|...)` commands with working defaults are welcome.
{{< /hint >}}

### S3 Storage Driver

A naive driver that treats the keys in an S3 cabaple storage as `/` delimited path names. While it does not support MOVE or etag propagation it can be used to read and write files. Better integration with native capabilities like versioning is possible but depends on the Use Case. Several storage solutions that provide an S3 interface also support some form of notifications that can be used to implement etag propagation.

## Data Providers

Clients using the CS3 API use an [InitiateFileDownload](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.InitiateFileDownloadRequest) and ]InitiateUpload](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.InitiateFileUploadRequest) request at the [reva gateway](https://cs3org.github.io/cs3apis/#cs3.gateway.v1beta1.GatewayAPI) to obtain a URL endpoint that can be used to eiter GET the file content or upload content using the resumable [tus.io](https://tus.io) protocol.

The *data provider* uses the same *storage driver* as the *storage provider* but can be scaled independently.

The dataprovider allows uploading the file to a quarantine area where further data analysis may happen before making the file accessible again. One use case for this is anti virus scanning for files coming from untrusted sources.

## Future work

### FUSE overlay filesystem
We are planning to further separate the concerns and use a local storage provider with a FUSE filesystem overlaying the actual POSIX storage that can be used to capture deletes and writes that might happen outside of ocis/reva.

It would allow us to extend the local storage driver with missing storage aspects while keeping a tree like filesystem that end users are used to see when sshing into the machine.

### Upload to Quarantine area
Antivirus scanning of random files uploaded from untrusted sources and executing metadata extraction or thumbnail generation should happen in a sandboxed system to prevent malicious users from gaining any information about the system. By spawning a new container with access to only the uploaded data we can further limit the attack surface.