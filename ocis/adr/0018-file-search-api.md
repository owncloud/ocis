---
title: "18. File Search API"
date: 2022-03-18T09:00:00+01:00
weight: 18
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0018-file-search-api.md
---

* Status: proposed
* Deciders: [@butonic](https://github.com/butonic), [@micbar](https://github.com/micbar), [@dragotin](https://github.com/dragotin), [@c0rby](https://github.com/c0rby)
* Date: 2022-03-18

## Context and Problem Statement

The ability to find files based on certain search terms is a key requirement for a system that provides the ability to store unstructured data on a large scale.

## Decision Drivers

* Have a simple yet powerful, scalable and performant way of finding files in oCIS
* Be able to construct intelligent searches based on metadata
* Allow the user to filter the search queries based on metadata

## Considered Options

* [Libre Graph API](#libre-graph-api)
* [WebDAV API](#webdav-api)

## Decision Outcome

Chosen option: [WebDAV API](#webdav-api) because the current WebUI is compatible with that API. We may use the GraphAPI later in a second iteration.

### Positive Consequences

* The existing Clients can continue to use the well-known API
* There are existing API tests which cover the basic behavior

### Negative consequences

* We have no server side result filtering capabilities

## Pros and Cons of the Options

### Libre Graph API

* Good, because we try to switch most of our HTTP requests to Libre Graph
* Good, because the Graph API supports scopes, sorting and query language
* Good, because it supports server side result filtering
* Bad, because there are currently no clients which support that

### WebDAV API

* Good, because WebDAV is a well-known and widely adopted Standard
* Good, because existing Clients continue to work without extra efforts
* Bad, because the syntax is limited
* Bad, because we cannot do server side result filtering

## Links

* [Search Indexing](0019-file-search-index.md)
* [Search Query Language](0020-file-search-query-language.md)
