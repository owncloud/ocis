# Search

The search service is responsible for metadata and content extraction, stores that data as index and makes it searchable. The following clarifies the extraction terms _metadata_ and _content_:

*   Metadata: all data that _describes_ the file like `Name`, `Size`, `MimeType`, `Tags` and `Mtime`.
*   Content: all data that _relates to content_ of the file like `words`, `geo data`, `exif data` etc.

## General Considerations

*   To use the search service, an event system needs to be configured for all services like NATS, which is shipped and preconfigured.
*   The search service consumes events and does not block other tasks.
*   When looking for content extraction, [Apache Tika - a content analysis toolkit](https://tika.apache.org) can be used but needs to be installed separately.

Extractions are stored as index via the search service. Consider that indexing requires adequate storage capacity - and the space requirement will grow. To avoid filling up the filesystem with the index and rendering Infinite Scale unusable, the index should reside on its own filesystem.

You can change the path to where search maintains its data in case the filesystem gets close to full and you need to relocate the data. Stop the service, move the data, reconfigure the path in the environment variable and restart the service.

When using content extraction, more resources and time are needed, because the content of the file needs to be analyzed. This is especially true for big and multiple concurrent files.

The search service runs out of the box with the shipped default `basic` configuration. No further configuration is needed, except when using content extraction.

Consider using a dedicated hardware for this service in case more resources are needed.

## Scaling

The search service can be scaled by running multiple instances. Some rules apply:

* With `SEARCH_ENGINE_BLEVE_SCALE=false`, which is the default , the search service has exclusive write access to the index. Once the first search process is started, any subsequent {search processes attempting to access the index are locked out.

* With `SEARCH_ENGINE_BLEVE_SCALE=true`, a search service will no longer have exclusive write access to the index. This setting must be enabled for all instances of the {search service.

## Search Engines

By default, the search service is shipped with [bleve](https://github.com/blevesearch/bleve) as its primary search engine. The available engines can be extended by implementing the [Engine](pkg/engine/engine.go) interface and making that engine available.

## Query language

By default, [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference) is used as query language, for an overview of how the syntax works, please read the [microsoft documentation](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference) for more details.

Not all parts are supported. The following list gives an overview of parts that are **not implemented** yet:

*   Synonym operators
*   Inclusion and exclusion operators
*   Dynamic ranking operator
*   ONEAR operator
*   NEAR operator
*   Date intervals

In the following [ADR](https://github.com/owncloud/ocis/blob/docs/ocis/adr/0020-file-search-query-language.md) you can read why we chose KQL.

## Extraction Engines

The search service provides the following extraction engines and their results are used as index for searching:

*   The embedded `basic` configuration provides metadata extraction which is always on.
*   The `tika` configuration, which _additionally_ provides content extraction, if installed and configured.

## Content Extraction

The search service is able to manage and retrieve many types of information. For this purpose the following content extractors are included:

### Basic Extractor

This extractor is the most simple one and just uses the resource information provided by Infinite Scale. It does not do any further analysis.

### Tika Extractor

This extractor is more advanced compared to the [Basic extractor](#basic-extractor). The main difference is that this extractor is able to search file contents. However, [Apache Tika](https://tika.apache.org/) is required for this task. Read the [Getting Started with Apache Tika](https://tika.apache.org/3.2.0/gettingstarted.html) guide on how to install and run Tika or use a ready to run [Tika container](https://hub.docker.com/r/apache/tika). See the [Tika container usage document](https://github.com/apache/tika-docker#usage) for a quickstart. Note that at the time of writing, containers are only available for the amd64 platform.

As soon as Tika is installed and accessible, the search service must be configured for the use with Tika. The following settings must be set:

*   `SEARCH_EXTRACTOR_TYPE=tika`
*   `SEARCH_EXTRACTOR_TIKA_TIKA_URL=http://YOUR-TIKA.URL`
*   `FRONTEND_FULL_TEXT_SEARCH_ENABLED=true`\
When using the Tika extractor, make sure to also set this enironment variable in the frontend service. This will tell the web client that full-text search has been enabled.

When the search service can reach Tika, it begins to read out the content on demand. Note that files must be downloaded during the process, which can lead to delays with larger documents.

Content extraction and handling the extracted content can be very resource intensive. Content extraction is therefore limited to files with a certain file size. The default limit is 20MB and can be configured using the `SEARCH_CONTENT_EXTRACTION_SIZE_LIMIT` variable.

When extracting content, you can specify whether [stop words](https://en.wikipedia.org/wiki/Stop_word) like `I`, `you`, `the` are ignored or not. Noramlly, these stop words are removed automatically. To keep them, the environment variable `SEARCH_EXTRACTOR_TIKA_CLEAN_STOP_WORDS` must be set to `false`.

When using the Tika container and docker-compose, consider the following:

*   See the [ocis_full](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_full) example.
*   Containers for the linked service are reachable at a hostname identical to the alias or the service name if no alias was specified.

## Search Functionality

The search service consists of two main parts which are file `indexing` and file `search`.

### Indexing

Every time a resource changes its state, a corresponding event is triggered. Based on the event, the search service processes the file and adds the result to its index. There are a few more steps between accepting the file and updating the index.

**IMPORTANT**

- When using the Tika Extractor, text and other data, such as EXIF data from images, are extracted from documents and written to the bleve index. Currently, this extra data cannot be searched. See the next section for more information.

### Search

A query via the search service will return results based on the index created.

**IMPORTANT**

- Though EXIF data can be present in the bleve index, currently, only text-related data can be extracted. Code must be written to make this type of extraction available to users.
- Currently, there is no ocis shell command or similar mechanism to view or browse the bleve index. This capability would be highly beneficial for developers and administrators to determine the type of data contained in the index.

### State Changes which Trigger Indexing

The following state changes in the life cycle of a file can trigger the creation of an index or an update:

#### Resource Trashed

The service checks its index to see if the file has been processed. If an index entry exists, the index will be marked as deleted. In consequence, the file won't appear in search requests anymore. The index entry stays intact and could be restored via [Resource Restored](#resource-restored).

#### Resource Deleted

The service checks its index to see if the file has been processed. If an index entry exists, the index will be finally deleted. In consequence, the file won't appear in search requests anymore.

#### Resource Restored

This step is the counterpart of [Resource Trashed](#resource-trashed). When a file is deleted, is isn't removed from the index, instead the service just marks it as deleted. This mark is removed when the file has been restored, and it shows up in search results again.

#### Resource Moved

This comes into play whenever a file or folder is renamed or moved. The search index then updates the resource location path or starts indexing if no index has been created so far for all items affected. See [Notes](#notes) for an example.

#### Folder Created

The creation of a folder always triggers indexing. The search service extracts all necessary information and stores it in the search index

#### File Created

This case is similar to [Folder created](#folder-created) with the difference that a file can contain far more valuable information. This gets interesting but time-consuming when data content needs to be analyzed and indexed. Content extraction is part of the search service if configured.

#### File Version Restored

Since Infinite Scale is capable of storing multiple versions of the same file, the search service also needs to take care of those versions. When a file version is restored, the service starts to extract all needed information, creates the index and makes the file discoverable.

#### Resource Tag Added

Whenever a resource gets a new tag, the service takes care of it and makes that resource discoverable by the tag.

#### Resource Tag Removed

This is the counterpart of [Resource tag added](#resource-tag-added). It takes care that a tag gets unassigned from the referenced resource.

#### File Uploaded - Synchronous

This case only triggers indexing if `async post processing` is disabled. If so, the service starts to extract all needed file information, stores it in the index and makes it discoverable.

#### File Uploaded - Asynchronous

This is exactly the same as [File uploaded - synchronous](#file-uploaded---synchronous) with the only difference that it is used for asynchronous uploads.

## Manually Trigger Re-Indexing a Space

The service includes a command-line interface to trigger re-indexing a space:

```shell
ocis search index --space $SPACE_ID
```

It can also be used to re-index all spaces:

```shell
ocis search index --all-spaces
```

Note that either `--space $SPACE_ID` or `--all-spaces` must be set.

## Notes

The indexing process tries to be self-healing in some situations.

In the following example, let's assume a file tree `foo/bar/baz` exists.
If the folder `bar` gets renamed to `new-bar`, the path to `baz` is no longer `foo/bar/baz` but `foo/new-bar/baz`.
The search service checks the change and either just updates the path in the index or creates a new index for all items affected if none was present.
