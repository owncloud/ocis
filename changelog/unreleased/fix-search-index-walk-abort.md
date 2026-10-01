Bugfix: Don't let unresolvable files abort a space reindex, and make the abort threshold configurable

`IndexSpace` aborts the walk after 5 consecutive files fail, on the assumption
that the content extractor is down. Two things made that fire on healthy
systems:

1. A file the walker listed but that could not be stat'ed by its path (a node
   whose stored name no longer matches its directory entry) counted as an
   extraction failure. The extractor was never involved.
2. The counter only ever saw files that actually needed work. On a space that
   is already almost fully indexed, five problem files scattered anywhere in
   walk order look "consecutive", the walk aborts at the same point every
   time, and nothing after that point is ever indexed until the files are
   fixed by hand.

Files that cannot be resolved are now logged and skipped without counting.
The threshold is configurable with `SEARCH_EXTRACTOR_MAX_CONSECUTIVE_FAILURES`
(default 5, unchanged); setting it to 0 never aborts, so a handful of files
the extractor cannot handle no longer keeps the rest of a space out of the
index.

https://github.com/owncloud/ocis/issues/13033
https://github.com/owncloud/ocis/pull/13034
