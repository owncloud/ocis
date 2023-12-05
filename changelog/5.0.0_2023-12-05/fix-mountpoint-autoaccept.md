Bugfix: set existing mountpoint on auto accept

When already having a share for a specific resource, auto accept would use custom mountpoints which lead to other errors. Now auto-accept is using the existing mountpoint of a share.

https://github.com/owncloud/ocis/pull/7592
