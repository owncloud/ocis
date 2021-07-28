---
title: "2. Persist accounts in a CS3 storage"
weight: 2
date: 2020-08-21T20:21:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0002-persist-accounts-using-cs3-storage.md
---

* Status: accepted
* Deciders: @butonic, @felixboehm
* Date: 2020-08-21

Technical Story: [File system based indexing](https://github.com/owncloud/ocis-accounts/pull/92)

## Context and Problem Statement

To set up High Availability (HA) or a geo-replicated setup we need to persist accounts in a distributed way. To efficiently query the accounts by email or username, and not only by id, they need to be indexed. Unfortunately, the [bleve](https://github.com/blevesearch/bleve) index we currently store locally on disk cannot be shared by multiple instances, preventing a scale out deployment.

## Considered Options

* Look into distributed bleve
* Persist users in a CS3 storage

## Decision Outcome

Chosen option: "Persist users in a CS3 storage", because we have one service less running and can rely on the filesystem for geo-replication and HA.

### Positive Consequences

* We can store accounts on the storage using the CS3 API, pushing geo-distribution to the storage layer.
* Backups of users and storage can be implemented without inconsistencies between using snapshots.

### Negative Consequences

* We need to spend time on implementing a reverse index based on files, and symlinks.
