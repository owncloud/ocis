Enhancement: Search content extraction

We've added the option to extract content while indexing resources.
To do so, the configuration flag `SEARCH_EXTRACTOR_TYPE` needs to be set to `tika` and `SEARCH_EXTRACTOR_TIKA_TIKA_URL` must point to a valid tika service url.

It's now also possible to provide different search engines, for the moment, bleve is the only supported engine.

https://github.com/owncloud/ocis/pull/4305
