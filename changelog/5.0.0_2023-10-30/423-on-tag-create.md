Bugfix: Return 423 status code on tag create

When a file is locked, return 423 status code instead 500 on tag create

https://github.com/owncloud/ocis/pull/7596
