---
title: "Decomposed FS on NFS"
date: 2020-03-15T16:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/storage-backends/
geekdocFilePath: dcfsnfs.md
---

{{< toc >}}

oCIS' default storage backend is the Decomposed FS. The Decomposed FS can be set up to run on a NFS share. That way the same storage can be provided over the network to other nodes running oCIS.

This document summarizes a few important considerations of the NFS setup and describes a tested setup. The test has not covered concurrent access of data from different nodes yet.

## NFS Server Setup

This document covers the linux kernel NFS server on a standard Linux running on x86_64.

The NFS server needs to be set up in a way that it supports [extended file attributes](https://en.wikipedia.org/wiki/Extended_file_attributes).

Extended attributes are supported by NFS starting with Kernel version 5.9, which means that the server with the NFS server has to run a kernel with that or a higher version number. To check that, run the command `uname -a` on the NFS server and compare the displayed version number.

The NFS server in the test setup was configured with the following line in the config file `/etc/exports`:

`/space/nfstest  192.168.178.0/24(rw,root_squash,async,subtree_check,anonuid=0,anongid=100,all_squash)`

This exports the directory `/space/nfstest` to the internal network with certain options.

Important:

- The share needs to be exported with the `async` option for proper NFS performance.

## NFS Client Mount

The nodes that run oCIS need to mount the NFS storage to a local mount point.

The test setup uses the client mount command: `mount -t nfs -o nfsvers=4 192.168.178.28:/space/nfstest /mnt/ocisdata/`

It sets the NFS version to 4, which is important to support extended attributes.

After successfully mounting the storage on the client, it can be checked if the NFS setup really supports extended attributes properly using the following commands.

`setfattr -n user.test -v "xattr test string" ocisdata/foo` to write an extended attribute to a file, and `getfattr -d ocisdata/foo` to list all the attributes a file has set.

{{< hint info >}}
The NFS server setup can be optimized considering system administrative-, performance- and security options. This is not (yet) covered in this documentation.
{{< /hint >}}

## oCIS Start using the NFS Share

The oCIS server can be instructed to set up the decomposed FS at a certain path by setting the environment variable `OCIS_BASE_DATA_PATH`.

The test setup started an oCIS tech preview single binary release using this start command:

```bash
./ocis init
OCIS_BASE_DATA_PATH=/mnt/ocisdata/ OCIS_LOG_LEVEL=debug OCIS_INSECURE=true PROXY_HTTP_ADDR=0.0.0.0:9200 OCIS_URL=https://hostname:9200  ./ocis server
```

This starts oCIS and a decomposed FS skeleton file system structure is set up on the NFS share.

The oCIS instance is passing a smoke test.
