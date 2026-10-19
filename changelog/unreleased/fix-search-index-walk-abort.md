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

The walk still resolves each file by its path. When Stat by path does not find
a file, the walk tries again with the ID it listed, so such a node is indexed.
A healthy file costs no extra GetPath call. A file that Stat does not find (CODE_NOT_FOUND, deleted
since it was listed) is skipped without counting; auth, gateway, transport and
deadline failures still count. Because the gateway also answers NOT_FOUND when
it cannot reach a storage provider, a walk in which at least 5 files, and more
than half of the files that needed work, were not found ends with an error
instead of reporting the space as indexed. Already-indexed, unchanged files
reset the counter, a cancelled or timed-out walk (including its retry waits)
stops with an error, and each walk logs a summary of indexed, unchanged, failed
and not-found files. The threshold became configurable with
`SEARCH_EXTRACTOR_MAX_CONSECUTIVE_FAILURES` (default 5, unchanged; 0 never
aborts and logs a warning at startup; negative values are rejected).

https://github.com/owncloud/ocis/pull/13034
https://github.com/owncloud/ocis/issues/13033
