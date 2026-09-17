Enhancement: Retry and abort on repeated extraction failures during indexing

During `ocis search index` bulk reindexing, if the content extractor (e.g.
Tika) becomes unavailable, individual file extractions are now retried up to
5 times with a 1-second delay between attempts. If a file still fails after
all retries, the failure is logged and the walk continues.

If 5 consecutive files fail extraction (indicating the extraction service is
down rather than a single file being problematic), the index walk is aborted
with an error so the admin can investigate.

Previously, extraction failures were silently logged and the walk continued
at full speed, accumulating thousands of wasted "connection refused" errors
when Tika was down.

https://github.com/owncloud/ocis/pull/12111
