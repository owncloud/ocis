Enhancement: Add a --timeout flag to the search index command

`ocis search index` called the search service with a fixed 10 minute request
timeout. The service walks the whole space inside that request, so once the
timeout expired the walk was cancelled with "error walking the tree: context
canceled" and everything after that point stayed unindexed. A space whose
walk takes longer than ten minutes, for example because the content extractor
is slow on a few files, could never be fully indexed from the command line.

The timeout is now a flag. `--timeout` takes a duration (default `10m`, so the
behaviour is unchanged out of the box); `0` waits indefinitely.

https://github.com/owncloud/ocis/issues/13052
