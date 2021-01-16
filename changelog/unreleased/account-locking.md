Enhancement: remove locking from accounts service

Tags: accounts, ocis-pkg

In the past we locked every request in the accounts service. This is problematic as a larger number of concurrent requests arrives at the accounts service.
The locking is now removed from the accounts service and is moved to the indexer.

Instead of doing locking for reads and writes we now differentiate them by using a named RWLock. 

- remove locking from accounts service
- add sync package with named mutex
- add named locking to indexer
- move cache into sync pkg

https://github.com/owncloud/ocis/pull/1212
https://github.com/owncloud/ocis/issues/966
