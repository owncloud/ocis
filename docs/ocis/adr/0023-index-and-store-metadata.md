---
title: "23. Index and store metadata"
date: 2023-10-17T15:15:00+01:00
weight: 23
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0023-index-and-store-metadata.md
---


* Status: accepted
* Deciders: [@butonic](https://github.com/butonic), [@theonering](https://github.com/dschmidt), [@micbar](https://github.com/micbar), [@dragotin](https://github.com/dragotin)
* Date: 2023-10-17

## Context and Problem Statement

ownCloud Infinite Scale is supposed to become a data platform and as such it needs to provide access to metadata.
Currently only metadata common to all file types (filesize, mime-type, ...) is stored in the index and the metadata storage.
We want to make other file type specific metadata available to consumers of our internal and external APIs.
Simple examples would be audio metadata like artist, album and title or exif metadata in images.

## Decision Drivers <!-- optional -->

## Considered Options

* [Store subset of extracted metadata required for graph api](#store-subset-of-extracted-metadata-required-for-graph-api)
* [Store subset of extracted metadata specified by another standard](#store-subset-of-extracted-metadata-specified-by-another-standard)
* [Store everything from extractors](#store-everything-from-extractors)

## Decision Outcome

Chosen option: "[store only subset of extracted metadata required for graph api](#store-subset-of-extracted-metadata-required-for-graph-api)", because Graph API is a simple common denominator and we want to avoid putting the complexity of mapping non-standardized data from potentially different extractors in several areas of the code base. Storage and index keys are determined by facet and property name, e.g. `audio.artist` for the artist in a music file. Storage keys are additionally prefixed with `libre.graph.`, i.e. `libre.graph.audio.artist`.
Handling Graph API specific metadata is a first step towards handling metadata. More generic and extensible handling of arbitrary metadata can be added later.

### Positive Consequences:

* Graph API endpoint implementation is trivial
* Documented public api and stored data are the same
* Reasonable complexity for the initial implementation

### Negative Consequences:

* Graph API is limited, so not *all* available metadata can be accessed
* Switching the internal format and adding more metadata later will require re-indexing

## Pros and Cons of the Options <!-- optional -->

### Store Subset of Extracted Metadata Required for Graph API

Use Graph API facets and properties for determining the subset of stored metadata and the storage key.
The index key for the `artist` property of the `audio` facet is `audio.artist`, the storage key is additionally prefixed with `libre.graph.`.

* Good, because central mapping of values happens consistently and only once in a central place
    - it happens in the extractor (integration) which likely knows best how to map metadata to standard properties
* Good, because when multiple extractors share a common set of provided values, applications can rely on the mapping and the complexity is kept low
* Bad, because not all metadata is available, not everything can be searched
* Good, because Graph API already chose a reasonable subset of most interesting properties

### Store Subset of Extracted Metadata Specified by Another Standard

There are a bunch of metadata standards but none of them is really universal. There is always something that is only supported in one or the other standard. Tika for example extracts audio metadata using a mixture of Dublin Core and XMP Dynamic Media keys.

- Bad, because it makes implementing a new extractor integration harder
- Bad, because it makes using the stored data more complicated than a simple standard like discussed above

### Store Everything from Extractors

- Good, because all metadata is available and searchable
- Good, because consuming applications can decide how to map data
- Good, because extractor implementation becomes more trivial
- Bad, because all applications become dependent on the extractor and need to handle different extractors on their own

## Links <!-- optional -->

* https://github.com/owncloud/libre-graph-api/pull/120 / https://learn.microsoft.com/de-de/graph/api/resources/audio?view=graph-rest-1.0
* https://github.com/owncloud/libre-graph-api/pull/122 / https://learn.microsoft.com/en-us/graph/api/resources/photo?view=graph-rest-1.0
* https://github.com/owncloud/libre-graph-api/pull/123 / https://learn.microsoft.com/en-us/graph/api/resources/geoCoordinates?view=graph-rest-1.0
* https://developer.adobe.com/xmp/docs/XMPNamespaces/xmpDM/
* https://www.dublincore.org/schemas/
