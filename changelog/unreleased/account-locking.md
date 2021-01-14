Enhancement: remove locking from accounts service

Tags: ocis

In the past we locked every request in the accounts service. This is problematic as soon as we start to hammer the system with many users at the same time.
The locking is now removed from the accounts service and is moved to the indexer.

Instead of doing locking for reads and writes we now differentiate them by using a named RWLock. 

- remove locking from accounts service
- add a cached named rwlock pkg
- use sync.map in the cache pkg
- use named rwlock in indexer pkg
- use sync.map in indexer pkg

https://github.com/owncloud/ocis/issues/966
