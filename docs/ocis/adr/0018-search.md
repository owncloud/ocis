
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

From the users perspective, the interface to search is just a single entry field where the user enters one or more search terms. The minimum expectation is that the search returns file names and links to files that

- have a file name that contains at least one of the search terms
- contain at least one of the search terms in the file contents
- have meta data that is equal or contains one of the search terms

More sophisticated search capabilities are expected and can be implemted, especially based on metadata.

Another assumption that this ADR makes is that the search operation is scoped to the file space. Each file space has its own search index, and the search query can be run in parallel per space.

## Decision Drivers

- Have a simple yet powerful way of finding files in oCIS
- Be able to construct intelligent searches based on metadata
- Allow the user to filter the search queries based on metadata

## Considered Options Query Notation

This lists the considered options for the search query notations.

### 1. Graph API

The search query adopts the Graph API to run search queries.

The [Libre Graph API](https://github.com/owncloud/libre-graph-api) would be inspired by the
[Microsoft Graph API](https://developer.microsoft.com/en-us/graph). Specifically the part that is [described here](https://docs.microsoft.com/en-us/graph/api/driveitem-search?view=graph-rest-1.0&tabs=http)

### 2. Keyword Query Language (KQL)

Implement a search based on the [Keyword Query Language (KQL)](https://github.com/SharePoint/sp-dev-docs/blob/master/docs/general-development/keyword-query-language-kql-syntax-reference.md), adopted from Sharepoint.

### 3. Simplified Query

Implement a very simple search approach: Return all files which contain at least one of the keywords in their name, path, alias or selected metadata.

## Considered Options Result Listing

The search request returns the result listing as described in the specification for Options 1. and 2.

For option 3. Simplified Query the result is returned in a ordered list of file Ids and relative paths that match the search pattern. Note that the file IDs are all within one file space.

## Considered Options Indexing

To start the indexing of a file the search service listens to create, update and delete events on the internal event bus of oCIS.

The events need to contain a valid reference that defines the file space and file id of the file in question. The event only must be sent when the file operation (update, creation, removal) is finished.

### Setting dirty Flags

*To be discussed*

Once a file is changed, it would be beneficial to set a metadata flag on the file that inicates that the file was changed and operations might have to happen, ie. propagating the index and updating the search index.

There should be a flag for every operation that is needed, ie. `user.dirty.etagpropagation` and `user.dirty.nameindex` as name, and the new ETag of the node as value.

The flags are set by each storage driver directly after the write was finished, within the write lock. The list of dirty flags to set needs to be pulled from a central method that lists them for all storage drivers.

### Multiple Indexes

*To be discussed*

For each space, it should be possible to have multiple indices.

For now we can foresee at least an index for file- and path segment names, and another for "simple" file meta data like times and permissions.

Benefit of having multiple indexing:
- Indexing can happen in parallel
- Querying can happen in parallel

## Decision Outcome

A search service is implemented as a separate microservice in oCIS. The service provides an API for search queries that deliver a list of results.

The indexing (create, update or remove) is only triggered asynchronously via events through the event bus. For that, the storage drivers need to send out the signals accordingly.
The search service provides an API to provide the following functionality:

### Query/Read

The search service provides a synchronous API to
- answer a search query based on search terms and storage space ids. The reply is a list of nodes that provides references to files within the space.

The search service API operates on one storage space by default. If a list of storage spaces to be searched through is provided as API parameter, the search is going sequentially through the list of Storage Spaces.

### Implementation for oCIS 2.0 GA

For oCIS 2.0, the search service only supports the Query Notation 3. Simplified Query.

It supports multiple indexes per space, but handles the query of them transparently to the outside. That means that the caller can not "choose" an index that should be queried or such.

The index is blieve based and saved as files via the CS3 API to the oCIS storage.

### Positive Consequences

### Negative Consequences

### Open Topics
