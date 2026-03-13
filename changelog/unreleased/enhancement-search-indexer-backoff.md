Enhancement: Add backoff and abort logic to search indexer when Tika is unavailable

The bulk indexer (IndexSpace) now detects repeated extraction failures
and pauses with exponential backoff instead of continuing to send
requests to an unreachable Tika server. After 5 consecutive failures
it pauses for 30 seconds, checks if Tika has recovered, and resumes.
After 5 backoff cycles with no recovery, it aborts the walk with a
clear error message. A summary of extracted, skipped, and failed
files is logged at the end of every index walk.

https://github.com/owncloud/ocis/pull/XXXXX
