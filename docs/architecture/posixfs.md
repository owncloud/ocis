---
title: "Posix Filesystem"
date: 2024-05-27T14:31:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/architecture
geekdocFilePath: posixfs.md
---

Posix FS is the working name for the collaborative storage driver for Infinite Scale.

The scope of this document is to give an high level overview to the technical aspects of the Posix FS and guide the setup.

## A Clean File Tree

Posix FS is a backend component that manages files on the server utilizing a "real" file tree that represents the data with folders and files in the file system as users are used to it. That is the big difference compared to Decomposed FS which is the default storage driver in Infinite Scale.

This does not mean that Infinte Scale is trading any of it's benefits to this new feature: It still implements simplicity by running without a database, it continues to store metadata in the file system and adds them transparently to chaches and search index, and it also features the full spaces concept as before, just to name a few example.

Based on the great system architecture of Infinite Scale, which allows to add different storage drivers with specific attributes, the Posix FS shares a lot of code with the decompsed FS which contributes to the stability and uniformity of both components.

However, the clearance of the file structure in the underlying file system is not the only benefit of the Posix FS. Moreover, this new technology allows users to manipulate the data directly in the file system, and any changes to files made even outside of Infinite Scale are monitored and directly reflected in Infinite Scale.

The first time ever with feature rich open source file synce & share systems, users can either choose to work with their data through the clients of the system, it's API's or directly in the underlying file system on the server.

That is another powerful vector for integration and enables a new universe of use cases across all domains. Just imagine how many software can write files, and can now directly make them accessible real time in a convenient, secure and efficient way.

## Technical Aspects

The PosixFS technology uses a few features of the underlying file system, which are directly contributing to the performance of the system.

While the simplest form of Posix FS runs with default file systems of every modern Linux system, the full power of this unfolds with more capable file systems such as IBM Storage Scale or Ceph. These are recommended as reliable foundations for big installations of Infinite Scale.

This chapter describes some technical aspects of the storage driver.

### Metadata

Infinite Scale is highly dependent on the efficient usage of meta data which are attached to file resources, but also logical elements such as spaces.

Metadata are stored in extended file attributes or in message pack files, as it is with other storage drivers, namely decompsed FS. All indexing and caching of metadata is located in higher system levels than the storage driver, and thus are not different to the components used with other storage drivers like the decomposed FS.

### Monitoring

To get information about changes that are done directly on the file system, Infinte Sale uses a monitoring system to directly watch the file system. This starts with the Linux inotify system and ranges to much more sophisticated services as for example in Spectrum Scale.

Based on the information transmitted by the watching service, Infinite Scale is able to "register" new or changed files into its own caches and internal management structures. That entitles Infinte Scale to deliver resource changes through the "traditional" channels such as APIs and clients.

Since the most important metadata is the file tree structure itself, the "split brain" situation between data and metadata is impossible to cause trouble.

### Automatic ETag Propagation

The ETag of a resource can be understood as a content fingerprint of any file- or folder resource in Infinite Scale. It is mainly used by clients to detect changes of resources. The rule is that if the content of a file changed the ETag has to change as well, as well as the ETag of all parent folders up to the root of the space.

A sophisticated underlying file system provides any attribute that fulfills this requirement and changes whenever content or metadata of a resource changes, and - which is most important - also changes the attribute of the parent resource and the parent of the parent etc.

If that is not available, Infinite Sale uses a built in mechanism to the maintain the ETag for each resource in the file meta data, and also propagates it automatically.

### Automatic Tree Size Propagation

Similar to the ETag propagation described before, Infinite Scale also tracks the accumulated tree size in all nodes of the file tree. A change to any file requires a re-calculation of the size attribute in all parent folders.

If the file system supports that natively that is a huge benefit.

### File Id Resolution

Infinite Scale uses an Id based approach to work with resources, rather than a file path based mechanism. The reason for that is that Id based lookups can be done way more efficient compared to tree traversals, just to name one reason.

The most important component of the Id is a unique file Id that identifies the resource within a space. Typically the Inode of a file could be used here. However, some file systems re-use inodes which must be avoided. Infinite Scale does not use the file Inode, but generates a UUID by default.

A powerful underlying file system would support Infinite Scale big times by providing an API that

1. Provides the Id for a given file path referenced resource
2. Provides the path for a given Id.

These two operations are very crucial for the performance of the entire system. For file systems that do not provide these APIs, Infinite Scale  provides internal caches to support the look ups.

### User Management

With the requirement that data can be manipulated either through the filesystem or the Infinite Scale system, the question under which uid the manipulation happens is an important question.

There are a few possible ways for user management:
1. Changes can either be only accepted by the same  user that Infinite Scale is running under, for example the user `ocis`. All manipulations in the filesystem have to be done by, and only by, this user.
2. Group based: All users who should be able to manipulate files have to be in a unix group. The Infinite Scale user has also to be in there. The default umask in the directory used has to allow group writing all over the place.
3. Impersonation: Infinite Scale impersonates the user who owns the folder on the file system to mimic the access as the user.

All possibilities have pros and cons for operations.

One for all, it seems reasonable to use LDAP to manage users which is the base for the Infinite Scale IDP as well as the system login system via PAM.

## Advanced Features

Depending on the capabilities of the underlying file system, the Infinite Scale PosixFS can benefit from more advanced funcitonality described here.

### Versioning

If the underlying file system is able to create versions of single resources (imagine a git based file system) this functionality could directly be used by Infinite Scale.

In the current state of the PosixFS, versioning is not supported.

### Trashbin

If the underlying file system handles deleted files in a trash bin that allows restoring of previously removed files, this functionality could directly be used by Infinite Scale.

If not available it will follow the [the Free Desktop Trash specificaton](https://specifications.freedesktop.org/trash-spec/trashspec-latest.html).

In the current state of the PosixFS, trash bin is not supported.

## Limitations

As of Q2/2024 the PosixFS is in technical preview state which means that it is not officially supported.

The tech preview comes with the following limitations:

1. User Management: Manipulations in the file system have to be done by the same user that runs Infinte Scale
2. Only inotify and GPFS file system change notification methods are supported
3. Advanced features versioning and trashbin are not supported yet
4. The space/project folders in the filesystem are named after the UUID, not the real space name

## Setup

This describes the steps to use the PosixFS storage driver with Infinite Scale.

It is possible to use different storage drivers in the same Infinite Scale installation. For example it is possible to set up one space running on PosixFS while others run decomposedFS.

### Prerequisites

To run PosixFS, the following prerequisites have to be fulfilled:

1. There must be storage available to store meta data and blobs, available under a root path
2. When using inotify, the storage must be local on the same machine. Network mounts do not work with Inotify
3. The storage root path must be writeable and executable by the same user Infinite Scale is running under
4. An appropiate version of Infinte Scale is installed, version number 5.0.5 and later
5. Either redis or nats-js-kv cache service


### Setup Configuration

This is an example configuration with environment variables that configures Infinite Scale to use PosixFS for all spaces it works with, ie. Personal and Project Spaces:

```
    "STORAGE_USERS_DRIVER": "posix",
    "STORAGE_USERS_POSIX_ROOT" : "/home/kf/tmp/posix-storage",
    "STORAGE_USERS_POSIX_WATCH_TYPE" : "inotifywait",
    "STORAGE_USERS_ID_CACHE_STORE": "nats-js-kv",              // for redis "redis"
    "STORAGE_USERS_ID_CACHE_STORE_NODES": "localhost:9233",    // for redis "127.0.0.1:6379"
```


