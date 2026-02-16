Bugfix: Do not reset received share state to pending

We fixed a problem where the states of received shares were reset to PENDING
in the "ocis migrate rebuild-jsoncs3-indexes" command

https://github.com/owncloud/ocis/issues/7319
