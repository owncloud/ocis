Bugfix: Prevent copying a file to a parent folder

When copying a file to a parent folder, the file would be copied to the parent folder, but the file would not be removed from the original folder.

https://github.com/owncloud/ocis/pull/8649
https://github.com/owncloud/ocis/issues/1230
https://github.com/cs3org/reva/pull/4571
`
