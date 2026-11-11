Bugfix: Prevent incomplete Tika extractions from permanently blocking re-index

When Tika returned HTTP 200 but its child processes (OCR, ImageMagick)
failed due to resource limits, the search index received metadata but
no content. The document was written to Bleve with the correct mtime,
and subsequent reindexes skipped it because the id+mtime check passed.
This left files permanently stuck as "indexed" with no searchable
content.

Two fixes are applied:

1. Validate Tika responses: if `MetaRecursive()` returns an empty
   metadata list, it is now treated as an extraction error so the
   document is not written to the index.

2. Add an `Extracted` field to indexed resources. It is set to `true`
   only after successful extraction. The reindex skip check now requires
   `Extracted:true`, so incompletely indexed documents are automatically
   re-processed on the next reindex run.

Note: existing search indexes will trigger a full re-extraction on the
next reindex because documents written before this change lack the
`Extracted` field.

https://github.com/owncloud/ocis/pull/12095
https://github.com/owncloud/ocis/issues/12093
