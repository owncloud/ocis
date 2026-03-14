Enhancement: Add backoff and abort when extraction fails during indexing

The bulk indexer (IndexSpace) now detects repeated extraction failures
and applies exponential backoff instead of continuing to send requests
to an unreachable extraction service. After 5 consecutive failures it
pauses for 30 seconds, then doubles the pause on each retry up to a
cap of 2 minutes. It only aborts after 30 minutes of continuous
failure. Any successful extraction resets the backoff entirely.
A summary of extracted, skipped, and failed files is logged at the
end of every index walk.

https://github.com/owncloud/ocis/pull/12111
