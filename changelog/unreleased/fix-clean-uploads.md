Bugfix: Fix clean uploads command

When using --clean ongoing uploads would be purged but the nodes would not be
reverted. This is now fixed.

https://github.com/owncloud/ocis/pull/11693
