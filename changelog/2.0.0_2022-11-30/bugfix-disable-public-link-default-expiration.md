Bugfix: Disable default expiration for public links

The default expiration for public links was enabled in the capabilities without providing a (then required) default amount of days for clients to pick a reasonable expiration date upon link creation. This has been fixed by disabling the default expiration for public links in the capabilities. With this configuration clients will no longer set a default expiration date upon link creation.

https://github.com/owncloud/ocis/issues/4445
https://github.com/owncloud/ocis/pull/4475
