---
title: "Collaborative Storage"
date: 2023-11-09T12:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/storage-backends/
geekdocFilePath: collaborative-storage.md
---

{{< toc >}}

One of the envisioned design goals of Infinite Scale is to work with so called _collaborative storage_, which means that the file system it is running on is not exclusive for Infinite Scale, but can be manipulated in parallel through third party tools. Infinite Scale is expected to monitor the changes that happen independently and react in a consistent and user friendly way.

A real world example of that would be a third party "data producer" that submits data directly into a file system path, not going through Infinite Scale APIs.

This document outlines a few challenges and design concepts for collaborative storage. It is also the base "checklist" for custom storage provider implementations for certain storages, ie. for Ceph- or IBM Storage Scale which provide features that allow more sophisticated and efficient implementations of this goal.


# Storage Driver Components

This discusses a few components and sub functions of the storage driver that have relevance for the collaborative storage.

## Path Locations

What is called "the oCIS file system" is defined as the entire filetree underneath a special path in the local POSIX file system, which might either be a real local file system or a mounted net filesystem. It is expected that oCIS is the only consumer of that file tree, except what is expected behaviour with a collaborative file system, that adds and edits files in that tree.

Underneath the oCIS file system root, there is an collection of different folders containing oCIS specific data.  Specific storage driver data is in the directory `storage/users`, organized by spaces.
(TODO: Check again how different storage drivers work together without overwriting data of each other)

## Spaces

Infinite Scale provides spaces as an additional organization organizational unit for data. Each space is a separate entity with its own attributes such as access patterns and quota.

A storage driver has to model the separation of spaces and provide a list of spaces in general and also a list of spaces a user can access. Furthermore, it needs to be able to create different types of spaces (Home- or Project space).

On POSIX, each space could for example be mapped to it's own directory in a special spaces folder under the oCIS root folder.

## ID to Path Lookup

Infinite Scale uses file IDs to efficiently identify files within a file tree. The lookup from a given ID to a path within the oCIS file tree is a very basic function that more or less defines the Infinite Scale performance. The functionality to for example query the file path for a given Inode number (which is the nearest equivalent for the Infinite Scale file ID) can not be done with standard POSIX system calls.

The interface defining the collaborative storage needs an abstraction for this particular function, returning the file id for a given path, and returning the path for given id.

## Change Notification

When a file is changed by a process outside of oCIS, this needs to be monitored by oCIS to quickly maintain internal caches and data structures as required.

The collaborative storage driver needs a way to achieve that. The easiest way for an POSIX based collaborative storage is inotify, that needs to be set up recursively on a file tree to record changes. Additional it is a challenge to destinguish between changes that were done from external activity and the ones that oCIS creates by its own file operations.

For GPFS, there is a subsystem called delivering that:

https://www.ibm.com/docs/en/storage-scale/5.1.9?topic=reference-clustered-watch-folder]

## ETag Propagation

ownCloud requires that changes which happen "down" in a tree, can be detected in the root element of the tree. That happens through the change of the ETag metadata of each file and/or directory. An ETag is a random, text based tag, that only has one requirement: It has to change its content if a resource further down in the file tree has changed either its content or its metadata. (See [this issue](https://github.com/owncloud/ocis/issues/3782) for further discussion about the ETag/CTag).

POSIX file systems do not maintain a change flag like the ETag by default. The file time stamps (atime, ctime, mtime) in general are not fine granular enough (only seconds for some file systems) and depend on the server time, which renders them useless in a distributed environment.

Infinite Scale needs to implement ETag propagation "up". For the collaborative storage, that needs to be combined with the change notification described above.

Certain file systems implement this functionality either independently from Infinite Scale (EOS) or at least support proper change notifications (Ceph, GPFS?).

## Metadata Management

Metadata are data "snippets" that are as tightly attached to files as ever possible. In best case, a rename of a file silently keeps the metadata as well. In POSIX, this can be achieved by extended file attributes with certain limitations.

## Quota

Each space has it's own quota, thus a storage driver implementation needs to consider that.

For GPFS for example, there is support for quota handling in the file system.

https://www.ibm.com/docs/en/gpfs/4.1.0.4?topic=interfaces-gpfs-quotactl-subroutine

Other systems store quota data in the metadata storage and implement propagation of used quota similar to the ETag propagation.

## User Management

With user management it is meant how to handle the users and groups within oCIS and how that reflects to the file system where data is stored.

### Exclusive Environment

In exclusive environments (aka. decomposedFS) all files of oCIS (ie. the entire oCIS filetree) belongs to a system user with the name `ocis` typically.

### Collaborative Storage

For collaborative storages, the approach described above does not longer work because users are supposed to be able to manipulate data in "their" file tree parts, and that is identified by ACLs and the owner of the files.

That requires a few prerequisites that have to be fulfilled:

1. oCIS as one "client" changing data and the system that allows to access the file tree directly have to use the same user provider, to ensure that each user that is available on a shell is also available in oCIS. That ensures that changes are authenticated through system ACLs and users. LDAP based authentication on the system via PAM and the same LDAP as source for the oCIS IDP should be a sufficient setup.
2. oCIS must be able to write as a "different" user than the ocis system user. That means that we somehow have to impersonate file changing ooperations and run these as the user that is authenticated in oCIS.

Example: There is a user ben. It has to have an entry in the LDAP that is used by IDP which oCIS is running "behind". With that, ben is able to authenticate through the IDP and work in the oCIS web app. The oCIS linux process will do writes and other changes impersonated as user ben.

For the access of data on the commandline, the logins to the linux system must be authenticated against the same LDAP - so that ben can authenticate on a terminal using username and password. With that, the user can interactively change data that belongs to user ben (simplified said).

To give permissions to groups, the linux group management must work accordingly. The same is true for file permissions.

## Trashbin

When a user deletes a file in oCIS it is moved to a so called trashbin that allows to restore the file if the deletion was accidentailly.

## Versions

When an existing file is changed, the former file state is to be preserved with data and metadata by oCIS. Some file system types provide this functionality via snapshots on partition or even file level. Other do not and have to implement that via a hidden directory keeping old file versions.




