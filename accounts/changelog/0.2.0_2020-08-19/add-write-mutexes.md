Bugfix: Add write mutexes

Concurrent account or groups writes would corrupt the json file on disk, because the different goroutines would be treated as a single thread from the os. We introduce a mutex for account and group file writes each. This locks the update frequency for all accounts/groups and could be further improved by using a concurrent map of mutexes with a mutex per account / group. PR welcome.

https://github.com/owncloud/ocis/accounts/pull/71
