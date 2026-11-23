Enhancement: Leave files whose extraction failed alone until they change

Every `IndexSpace` walk sent every file that was not indexed to the content
extractor again, including files whose extraction had already failed. A file
that makes Tika time out (a large scanned PDF going through OCR, for example)
therefore cost a full extractor timeout and a restart of Tika's forked worker on
every walk, and every upload, move or restore in the space starts a walk. On one
installation four such PDFs caused 122 Tika worker restarts in two days and kept
the extractor busy for about a third of the hours the space was in use. A file
whose extraction failed was not indexed at all, so it could not be found by its
name either.

With the new `SEARCH_EXTRACTOR_FAILED_RETRY_INTERVAL` set (for example `24h`), a
failed extraction during a space walk is recorded with the time of the failure.
A file with no index entry yet is indexed by its name and metadata only; a file
that already has one keeps its indexed content and metadata. Later walks skip
the file while it is unchanged and the interval has not passed. Any change of
its mtime or size, restoring an older version included, has it extracted again
straight away. Only failures that are down to the file are recorded: not a
failure on upload (the next walk tries the file again), not one where every
attempt failed to reach the extractor, not one for a file that was deleted in
the meantime, and the failures of a run that aborts the walk on
`SEARCH_EXTRACTOR_MAX_CONSECUTIVE_FAILURES` are forgotten again. Skipped files
neither count towards nor reset that threshold, and the walk summary reports
them as `failed_earlier`. The default, 0, keeps the current behaviour: every
walk retries such files and nothing is recorded.

The new index field is stored but not searchable in indexes created from this
version on. An existing index maps it dynamically, so its value (a timestamp)
is also searchable there; this is harmless.

https://github.com/owncloud/ocis/pull/13506
