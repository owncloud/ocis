
---
title: "8. oCIS Search Infrastructure"
date: 2021-06-08T09:00:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0008-search.md
---

* Status: accepted
* Deciders: @butonic, @micbar, @dragotin
* Date: 2021-06-08

Technical Story: oCIS Internal Services and APIs for File Search

## Context and Problem Statement

The ability to find files based on certain search terms is a key requirement for a system that provides the ability to store unstructured data large scale. This ADR outlines the concepts to implement search in oCIS.

From the users perspective, the interface to search is just a single entry field where the user enters one or more search terms. The expectation is that the search returns file names and links to it of files that

- have a file name that contains the search terms
- contain the search terms
- have meta data that is equal or contains the search terms

More sophisticated search capabilities are expected and can be implemted, especially based on metadata.

## Decision Drivers

- Have a simple yet powerful way of finding files in oCIS
- Be able to construct intelligent searches based on metadata
- Allow the user to filter the search queries based on metadata

## Considered Options

1. [Microsoft Graph API](https://developer.microsoft.com/en-us/graph) inspired API, specifically the part that is [described here](https://docs.microsoft.com/en-us/graph/api/driveitem-search?view=graph-rest-1.0&tabs=http)

2. Implement a search based on the [Keyword Query Language (KQL)](https://github.com/SharePoint/sp-dev-docs/blob/master/docs/general-development/keyword-query-language-kql-syntax-reference.md)

## Decision Outcome

A search service is implemented as a microservice in oCIS. This service provides an API to

- answer a search query based on a search term. The reply is a list of nodes that provide links within the storage space.
- Index a new or updated file based on it's mime type.

The search service API operates on one Storage Space by default. If a list of storage spaces to be searched through is provided as API parameter, the search is going sequentially through the list of Storage Spaces.

Each Storage Space has its very own and independant search index. That way, many search services can be started to search for files in all services in parallel. The search index is either file based within the storage space, or delegated into a dedicated service such as Elastic Search. The configuration which method to use is encapsulated in the search service.

As a first implementation, the API described in [the stable version of the Graph search API]((https://docs.microsoft.com/en-us/graph/api/driveitem-search?view=graph-rest-1.0&tabs=http) will be sufficient.

Later, the KQL will be implemented to make more use of stored metadata.

### Positive Consequences

### Negative Consequences


### Open Topics
