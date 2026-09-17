Bugfix: Stop logging every service re-registration

Services that publish themselves to the service registry re-register on a timer
to keep their TTL alive. Every one of those refreshes wrote a debug line
carrying the whole service record, so an idle deployment running at debug level
produced a steady stream of log entries with no user interaction at all. The
refresh is now silent unless it fails, which is the case that was worth a log.

https://github.com/owncloud/ocis/pull/12833
https://github.com/owncloud/ocis/issues/8274
