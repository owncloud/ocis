---
title: "24. Messagepack metadata"
date: 2024-02-09T14:57:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0024-msgpack-metadata.md
---


* Status: accepted
* Deciders: [@butonic](https://github.com/butonic), [@aduffeck](https://github.com/aduffeck), [@micbar](https://github.com/micbar), [@dragotin](https://github.com/dragotin)
* Date: [2023-03-15](https://github.com/cs3org/reva/pull/3711/commits/204253eee9dbb8e7fa93a01f3f94a2d28ce40a06)

## Context and Problem Statement

File metadata management is an important aspect for oCIS as a data platform. While using extended attributes to store metadata allows attaching the metadata to the actual file it causes a significant amount of syscalls that outweigh the benefits. Furthermore, filesystems are subject to different limitations in the number of extended attributes or the value size that is available.

## Decision Drivers <!-- optional -->

Performance of reading extended attributes suffers from the syscall overhead when listing and reading all attributes. Getting rid of limitations imposed by the filesystem used to store decomposedfs metadata. 

## Considered Options

Going back to the original [ADR-0016 Storage for Files Metadata]({{< ref "0016-files-metadata.md" >}}) we decided to use a dedicated file for metadata storage next to the decomposedfs file representing the node. Several options for the data format were considered:

* Use JSON files to store metadata
* Use INI files to store metadata
* Use msgpack files to store metadata
* Use protobuf messages to store metadata

## Decision Outcome

Chosen option: "[msgpack files](#msgpack-files)", because we want to stay with a self describing binary format. This is a performance tradeoff that is faster and more efficient than text based formats and more flexible but less efficient than protobuf.

Note: directory listings are still read from the storage and remain uncached.

### Positive Consequences:

* Way less syscalls
* Node metadata can easily be cached, avoiding all trips to the storage until a file changes.

### Negative Consequences:

* We need to migrate existing metadata
* We need to build tooling that allows manipulating metadata similar to `setfattr` and `getfattr`.

## Pros and Cons of the Options <!-- optional -->

### Ini files

* Good, human readable
* Good, self describing
* Good, widely used and well understood
* Good, suited for key value like content - exactly what we need for extended attributes
* Bad, slower and less efficient than binary formats

### JSON files

* Good, human readable
* Good, self describing
* Good, widely used and well understood
* Good, could be used for more than just key value
* Bad, slower and less efficient than binary formats

### Msgpack files

* Good, self describing
* Good, efficient because it is binary encoded
* Good, could be used for more than just key value
* Bad, not human readable - requires tooling to manipulate safely

### protobuf files

* Good, very efficient because it is binary encoded
* Good, could be used for more than just key value
* Bad, not human readable
* Bad, not self describing - requires tooling to evolve the messages

## Links <!-- optional -->

* supersedes [ADR-0016 Storage for Files Metadata]({{< ref "0016-files-metadata.md" >}})
* [The need for speed — Experimenting with message serialization](https://medium.com/@hugovs/the-need-for-speed-experimenting-with-message-serialization-93d7562b16e4)