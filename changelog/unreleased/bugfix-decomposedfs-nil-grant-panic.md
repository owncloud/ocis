Bugfix: Fix a crash when the storage reads permissions

The storage service could stop when it read the permissions of a user. We
updated the reva library. The new version checks the result of a read before
it uses the result. The new version also rejects invalid permission data and
does not stop.

https://github.com/owncloud/ocis/pull/13436
