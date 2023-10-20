Enhancement: Tika content extraction cleanup for search

So far it has not been possible to determine whether
the content for the search should be cleaned of stop words or not.

This can now be set with the newly introduced settings option `SEARCH_EXTRACTOR_TIKA_CLEAN_STOP_WORDS=false`
which is enabled by default.

In addition, the stop word cleanup is no longer as aggressive and now ignores numbers, urls,
basically everything except the defined stop words.

https://github.com/owncloud/ocis/pull/7553
https://github.com/owncloud/ocis/issues/6674
