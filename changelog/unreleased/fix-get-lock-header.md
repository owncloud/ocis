Bugfix: Return wopi lock header in get lock response

We fixed a bug where the wopi lock header was not returned in the get lock response. This is now fixed and the wopi validator tests are passing.

https://github.com/owncloud/ocis/pull/10469
