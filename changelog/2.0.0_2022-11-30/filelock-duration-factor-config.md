Bugfix: decomposedfs increase filelock duration factor

We made the file lock duration per lock cycle for decomposedfs configurable and increased it to make locks work on top of NFS.

https://github.com/owncloud/ocis/pull/5130
https://github.com/owncloud/ocis/issues/5024
