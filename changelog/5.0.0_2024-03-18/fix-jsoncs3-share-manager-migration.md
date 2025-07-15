Bugfix: Do not reset state of received shares when rebuilding the jsoncs3 index

We fixed a problem with the "ocis migrate rebuild-jsoncs3-indexes" command which reset the state of received shares to "pending".

https://github.com/owncloud/ocis/issues/7319
