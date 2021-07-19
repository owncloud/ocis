---
title: "4. Support Hot Migration"
weight: 4
date: 2020-12-09T20:21:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0004-support-hot-migration.md
---

* Status: proposed
* Deciders: @butonic, @micbar, @dragotin, @hodyroff, @pmaier1
* Date: 2021-03-16

Technical Story: \[description | ticket/issue URL\]

## Context and Problem Statement

Migration is one of the most important topics of the oCIS story. We need to provide a concept how to migrate from oC10 to oCIS.

## Decision Drivers

- Do not lose file blob or meta data.
  - To prevent a sync surge from clients the etag for files should be migrated.
  - To prevent internal links from breaking or pointing to wrong files the file id of existing files needs to be migrated.
  - To prevent data loss trash and version blobs should be migrated.
- Existing shares like public links and federated shares must remain functional after the migration.
  - To prevent internal shares the share type, permissions and expiry needs to be migrated.
  - To prevent public links from breaking the url token, permissions, expiry and password needs to be migrated.
  - *What about federated shares?*
  - *What about additional share permissions, eg. comment on office files?*
- Legacy clients need to keep working
  - To keep existing clients working the `remote.php/webdav` and `dav/files/<username>` webdav endpoints as well as the ocs API need to be available.
- *What about [app passwords/tokens](https://doc.owncloud.com/server/user_manual/personal_settings/security.html#app-passwords-tokens)?*

## Considered Options

1. Cold Migration: migrate data while systems are not online, so no user interaction happens in between.
2. Hot Migration: one or both systems are online during migration.

## Decision Outcome

Chosen option: "\[option 1\]", because \[justification. e.g., only option, which meets k.o. criterion decision driver | which resolves force force | … | comes out best (see below)\].

### Positive Consequences

- \[e.g., improvement of quality attribute satisfaction, follow-up decisions required, …\]
- …

### Negative Consequences

- \[e.g., compromising quality attribute, follow-up decisions required, …\]
- …

## Pros and Cons of the Options

### Cold Migration

The migration happens while the service is offline. File metadata, blobs and share data is exported from ownCloud 10 and imported in oCIS. This can happen user by user, where every user export would contain the file blobs, their metadata, trash, versions, shares and all metadata that belongs to the users storage. To prevent group shares from breaking, users in the same groups must be migrated in batch. Depending on the actual group shares in an instance this may effectively require a complete migration in a single batch. 

- Good, because oCIS can be tested in a staging system without writing to the production system.
- Good, because file layout on disk can be changed to support new storage driver capabilities.
- Bad, because the export and import might require significant amounts of storage.
- Bad, because a rollback to the state before the migration might cause data loss of the changes that happend in between.
- Bad, because the cold migration can mean significant downtime.

### Hot Migration

The migration happens in subsequent stages while the service is online.

- Good, because the admin can migrate users from old to new backend in a controlled way.
- Good, because users and admins can learn to trust the new system.
- Good, because there can be preparations even long before the migrations happens in parallel on the oC10 codebase, ie. addition of metadata that is needed while the system operates.
- Good, because the downtime of the system can be fairly small.
- Bad, because it is more complex and might drag on for a long time.


## Links

        
- [Clarify responsibilities of share providers and storage providers · Issue #1377 · cs3org/reva (github.com)](https://github.com/cs3org/reva/issues/1377) because the share manager for oCIS should store share information on the storage system. And [storageprovider should persist share creator · Issue #93 · cs3org/cs3apis (github.com)](https://github.com/cs3org/cs3apis/issues/93) finally: [eos: store share id in inherited xattr · Issue #543 · cs3org/reva (github.com)](https://github.com/cs3org/reva/issues/543)
