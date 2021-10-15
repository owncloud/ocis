Bugfix: update glauth to 20210729125545-b9aecdfcac31

* Fixes the backend config not being passed correctly in ocis
* Fixes a mutex being copied, leading to concurrent writes
* Fixes UTF8 chars in filters
* Fixes case insensitive strings

https://github.com/owncloud/ocis/pull/2336
https://github.com/glauth/glauth/pull/198
https://github.com/glauth/glauth/pull/194