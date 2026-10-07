Bugfix: Reindex no longer aborts on unresolvable files; threshold configurable

`IndexSpace` aborted the walk after 5 consecutive files failed, on the assumption
that the content extractor was down. Two things made that fire on healthy
systems:

1. A file the walker listed but that could not be stat'ed by its path (a node
   whose stored name no longer matched its directory entry) counted as an
   extraction failure. The extractor was never involved.
2. The counter only ever saw files that actually needed work. On a space that
   was already almost fully indexed, five problem files scattered anywhere in
   walk order looked "consecutive", the walk aborted at the same point every
   time, and nothing after that point was indexed until the files were fixed
   by hand.

Files that Stat no longer finds by their path (CODE_NOT_FOUND) were changed to
be logged and skipped without counting; auth, gateway, transport and deadline
failures still count. Already-indexed, unchanged files now reset the counter,
and a cancelled or timed-out walk stops with an error instead of skipping the
rest of the space. The threshold became configurable with
`SEARCH_EXTRACTOR_MAX_CONSECUTIVE_FAILURES` (default 5, unchanged; 0 never
aborts; negative values are rejected).

https://github.com/owncloud/ocis/issues/13033
https://github.com/owncloud/ocis/pull/13034
