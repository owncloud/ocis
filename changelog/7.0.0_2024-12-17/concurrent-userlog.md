Enhancement: Concurrent userlog processing

We now start multiple go routines that process events. The default of 5 goroutines can be changed with the new `USERLOG_MAX_CONCURRENCY` environment variable.

https://github.com/owncloud/ocis/pull/10504
