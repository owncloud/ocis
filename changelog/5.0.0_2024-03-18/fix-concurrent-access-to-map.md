Bugfix: Fix concurrent access to a map

We fixed the race condition that led to concurrent map access in a publicshare manager.

https://github.com/owncloud/ocis/pull/8269
https://github.com/cs3org/reva/pull/4472
https://github.com/owncloud/ocis/issues/8255
