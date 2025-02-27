---
title: "19. File Search Index"
date: 2022-03-18T09:00:00+01:00
weight: 19
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0019-file-search-index.md
---

* Status: accepted
* Deciders: [@butonic](https://github.com/butonic), [@micbar](https://github.com/micbar), [@dragotin](https://github.com/dragotin), [@c0rby](https://github.com/c0rby)
* Date: 2022-03-18

## Context and Problem Statement

The ability to find files based on certain search terms is a key requirement for a system that provides the ability to store unstructured data on a large scale.

More sophisticated search capabilities are expected and can be implemented, especially based on metadata.

To trigger the indexing of a file, the search service listens to create, update and delete events on the internal event bus of oCIS.

The events need to contain a valid reference that defines the file space and file id of the file in question. The event only must be sent when the file operation (update, creation, removal) is finished.

Sharing adds more complexity because the index also needs to react to create, delete and modify shares events. Sharing should not duplicate the indexed data, especially within spaces or group shares.

## Decision Drivers

* Have a simple yet powerful, scalable and performant way of finding files in oCIS
* Be able to construct intelligent searches based on metadata
* Allow the user to filter the search queries based on metadata
* Basic File Search needs to be implemented out of the box without external dependencies
* The Search Indexing Service should be replaceable with more sophisticated technologies like Elasticsearch
* Make use of the spaces architecture to shard search indexes by space
* The Search Indexing Service needs to deal with multiple users accessing the same resources due to shares
* The Search Service should be compatible with different search indexing technologies

## Considered Options

* [Bleve Search](#bleve-search)
* [Elastic Search](#elastic-search)

## Decision Outcome

Chosen option: Bleve Search, because we can fulfill the MVP and include it into the single binary.

### Positive Consequences

* Basic File Search works out of the box
* We do not need heavy external dependencies which need to be deployed alongside

### Negative consequences

* We need to be aware of the scaling limits
* We need to find a way to work with shares and spaces
* It has a limited query language

## Pros and Cons of the Options

### Bleve Search

* Good, because it is written in GoLang and can be bundled into the single oCIS binary
* Good, because it is a lightweight but powerful solution which could fulfill a lot of use cases
* Bad, because we do not know exactly how we can represent shares in the index without duplicating data
* Bad, because it is a single process
* Bad, because the query language is limited

### Elastic Search

* Good, because it has become an industry standard
* Good, because it supports a rich query language
* Good, because it has built in cluster support and scales well
* Good, because it has a permission system and supports multiple users and groups to access the same resource
* Bad, because it is a heavy setup and needs extra effort and knowledge

## Links

* [Search API](0018-file-search-api.md)
* [Search Query Language](0020-file-search-query-language.md)
* [Bleve Search on GitHub](https://github.com/blevesearch/bleve)
* [ElasticSearch](https://www.elastic.co/elastic-stack/)
