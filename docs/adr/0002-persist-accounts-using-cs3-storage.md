# 2. Persist accounts in a CS3 storage

* Status: accepted <!-- optional -->
* Deciders: @butonic, @felixboehm <!-- optional -->
* Date: 2020-08-21 <!-- optional -->

Technical Story: [File system based indexing](https://github.com/owncloud/ocis-accounts/pull/92) <!-- optional -->

## Context and Problem Statement

To set up HA or a geo replicated setup we need to persist accounts in a distributed way. Furthermore, the bleve index makes the accounts service stateful, which we wand to avoid for a scale out deployment.

## Considered Options

* Look into distributed bleve
* Persist users in a CS3 storage

## Decision Outcome

Chosen option: "Persist users in a CS3 storage", because we have one service less running and can rely on the filesystem for geo replication and HA.

### Positive Consequences <!-- optional -->

* We can store accounts on the storage using the CS3 API, pushing geo-distribution to the storage layer.
* Backups of users and storage can be implemented without inconsistencies between using snapshots.

### Negative Consequences <!-- optional -->

* We need to spend time on implementing a reverse index based on files, and symlinks.
