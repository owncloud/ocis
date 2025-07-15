---
title: "20. File Search Query Language"
date: 2022-06-23T09:00:00+01:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0020-file-search-query-language.md
---

* Status: accepted
* Deciders: [@butonic](https://github.com/butonic), [@micbar](https://github.com/micbar), [@dragotin](https://github.com/dragotin), [@c0rby](https://github.com/c0rby), [@kulmann](https://github.com/kulmann), [@felix-schwarz](https://github.com/felix-schwarz), [@JammingBen](https://github.com/JammingBen)
* Date: 2023-06-23

## Context and Problem Statement

From the users perspective, the interface to search is just a single form field where the user enters one or more search terms. The minimum expectation is that the search returns file names and links to files that:

* have a file name that contains at least one of the search terms
* contain at least one of the search terms in the file contents
* have metadata that is equal or contains one of the search terms

## Decision Drivers

* The standard user should not be bothered by a query syntax
* The power user should also be able to narrow his search with an efficient and flexible syntax
* We need to consider different backend technologies which we need to access through an abstraction layer
* Using different indexing systems should lead to a slightly different feature set without changing the syntax completely

## Considered Options

* [KQL - Keyword Query Language](#keyword-query-language)
* [Simple Query](#simplified-query)
* [Lucene Query Language](#lucene-query-language)
* [Solr Query Language](#solr-query-language)
* [Elasticsearch Query Language](#elasticsearch-query-language)

## Decision Outcome

Chosen option: [KQL - Keyword Query Language](#keyword-query-language), because it enables advanced search across all platforms.

### Positive Consequences

* We can use the same query language in all clients

### Negative consequences

* We need to build and maintain a backend connector

## Pros and Cons of the Options

### Keyword Query Language

The Keyword Query Language (KQL) is used by Microsoft Share Point and other Microsoft Services. It uses very simple query elements, property restrictions and operators.

* Good, because we can fulfill all our current needs
* Good, because it is very similar to the used query language in iOS
* Good, because it supports date time keywords like "today", "this week" and more
* Good, because it can be easily extended to use "shortcuts" for eg. document types like `:presentation` which combine multiple mime types.
* Good, because it is successfully implemented and used in similar use cases
* Good, because it gives our clients the freedom to always use the same query language across all platforms
* Good, because Microsoft Graph API is using it, we will have an easy transition in the future
* Bad, because we need to build and maintain a connector to different search backends (bleve, elasticsearch or others)

### Simplified Query

Implement a very simple search approach: Return all files which contain at least one of the keywords in their name, path, alias or selected metadata.

* Good, because that covers 80% of the users needs
* Good, because it is very straightforward
* Good, because it is a suitable solution for GA
* Bad, because it is below the industry standard
* Bad, because it only provides one search query

### Lucene Query Language

The Lucene Query Parser syntax supports advanced queries like term, phrase, wildcard, fuzzy search, proximity search, regular expressions, boosting, boolean operators and grouping. It is a well known query syntax used by the Apache Lucene Project. Popular Platforms like Wikipedia are using Lucene or Solr, which is the successor of Lucene

* Good, because it is a well documented and powerful syntax
* Good, because it is very close to the Elasticsearch and the Solr syntax which enhances compatibility
* Bad, because there is no powerful and well tested query parser for golang available
* Bad, because it adds complexity and fulfilling all the different query use-cases can be an "uphill battle"

### Solr Query Language

Solr is highly reliable, scalable and fault-tolerant, providing distributed indexing, replication and load-balanced querying, automated failover and recovery, centralized configuration and more. Solr powers the search and navigation features of many of the world's largest internet sites.

* Good, because it is a well documented and powerful syntax
* Good, because it is very close to the Elasticsearch and the Lucene syntax which enhances compatibility
* Good, because it has a strong community with large resources and knowledge
* Bad, because it adds complexity and fulfilling all the different query use-cases can be an "uphill battle"

### Elasticsearch Query Language

Elasticsearch provides a full Query DSL (Domain Specific Language) based on JSON to define queries. Think of the Query DSL as an AST (Abstract Syntax Tree) of queries, consisting of two types of clauses. It is able to combine multiple query types into compound queries. It is also a successor of Solr.

* Good, because it is a well documented and powerful syntax
* Good, because it is very close to the Elasticsearch and the Solr syntax which enhances compatibility
* Good, because there is a stable and well tested go client which brings a query builder
* Good, because it could be used as the query language which supports different search backends by just implementing what is needed for our use-case
* Bad, because it adds complexity and fulfilling all the different query use-cases can be an "uphill battle"

## Links

* [Search API](0018-file-search-api.md)
* [Search Indexing](0019-file-search-index.md)
* [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference)
* [Apache Lucene](https://lucene.apache.org/)
* [Apache Solr](https://solr.apache.org/)
* [Elastic Search](https://solr.apache.org/)
* [Elastic Search for go](https://github.com/elastic/go-elasticsearch)
