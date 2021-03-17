# 4. Support Hot Migration

* Status: proposed
* Deciders: @butonic, @micbar, @dragotin, @hodyroff, @pmaier1
* Date: 2021-03-16

Technical Story: \[description | ticket/issue URL\]

## Context and Problem Statement

Migration is one of the most important topics of the oCIS story. We need to provide a concept how to migrate from oC10 to oCIS.

## Decision Drivers

- Do not lose file blob or meta data.
- Existing shares like public links and federated shares must remain functional after the migration.
- Legacy clients need to keep working

## Considered Options

- Cold Migration
- Hot Migration

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

The migration happens while the service is offline. File metadata, blobs and share data is exported from ownCloud 10 and imported in oCIS.

- Good, because oCIS can be tested in a staging system without writing to the production system.
- Good, because file layout can be changed together with the architecture migration.
- Bad, because the export and import might require significant amounts of storage.
- Bad, because a rollback to the state before the migration might cause data loss of the changes that happend in between.

### Hot Migration

The migration happens in subsequent stages while the service is online.

- Good, because users can switch between the backends on the fly.
- Good, because users and admins can learn to trust the new system.
- Bad, because it is more complex and might drag on for a long time.


## Links

        
- [Clarify responsibilities of share providers and storage providers · Issue #1377 · cs3org/reva (github.com)](https://github.com/cs3org/reva/issues/1377) because the share manager for ocis should store share information on the storage system. And [storageprovider should persist share creator · Issue #93 · cs3org/cs3apis (github.com)](https://github.com/cs3org/cs3apis/issues/93) finally: [eos: store share id in inherited xattr · Issue #543 · cs3org/reva (github.com)](https://github.com/cs3org/reva/issues/543)