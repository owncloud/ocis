---
title: "25. Distributed Search Index"
date: 2024-02-09T16:27:00+01:00
weight: 25
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0025-distributed-search-index.md
---

* Status: draft
* Deciders: [@butonic](https://github.com/butonic), [@fschade](https://github.com/fschade), [@aduffeck](https://github.com/aduffeck)
* Date: 2024-02-09

## Context and Problem Statement

Search is currently implemented with [blevesearch](https://github.com/blevesearch/bleve), which internally uses bbolt. bbolt writes to a local file, which prevents scaling out the service.

The initial implementation used a single blevesearch index for all spaces. While this makes querying all spaces easy because the results do not need to be aggregated from multiple indexes, the single node becomes a bottleneck when answering search queries. Furthermore, indexing is also part of the search service and has to share the resources.

## Decision Drivers <!-- optional -->

* Indexing should be decoupled from the search service
* The search service should be able to scale horizontally
* The solution needs to be embeddable in the single binary

## Considered Options
 
* one index per space
* [elasticsearch](https://github.com/elastic/elasticsearch) (java)
* [dgraph](https://github.com/dgraph-io/dgraph) (go)
* [manticore](https://github.com/manticoresoftware/manticoresearch/) (C++)
* [meilisearch](https://github.com/meilisearch/meilisearch) (Rust) 

## Decision Outcome

Chosen option: *???*

### Positive Consequences:

* TODO

### Negative Consequences:

* TODO

## Pros and Cons of the Options <!-- optional -->

### one index per space

Instead of using a single index (current implementation) or a distributed search index like elasticsearch the search service should aggregate queries from dedicated indexes per space. The api to a space index provider should be able to take multiple space ids in the request, similar to how a storage provider can handle multiple spaces. When treating spaces and the corresponding search index to belong together we can also treat them as a single unit for backup and restore. In federated deployments we can send the search queries to all search providers / spaces that the user has access to.

How a search provider is implemented then depends on the requirements. For a single node deployment bleve might be fine, for a kubernetes deployment a dedicated service might be the better fit.

### elasticsearch

* Good, commercial support available at https://www.elastic.co/de/pricing
* Good, industry standard
* Bad, nobody seems to like it
* Bad, not embeddable (Java)

### dgraph

* Good, commercial support available at https://dgraph.io/pricing
* Good, embeddable? (go) - TODO verify

### manticore
* Good, commercial support available at https://manticoresearch.com/services/
* Bad, not embeddable (C++)

### meilisearch
* Good, commercial support available at https://www.meilisearch.com/pricing
* Bad, not embeddable (Rust)

## Links <!-- optional -->

* supersedes [ADR-0019 File Search Index]({{< ref "0019-file-search-index.md" >}})