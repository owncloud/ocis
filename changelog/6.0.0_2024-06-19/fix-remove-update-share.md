Bugfix: Fix remove/update share permissions

This is a workaround that should prevent removing or changing the share permissions when the file is locked.
These limitations have to be removed after the wopi server will be able to unlock the file properly.
These limitations are not spread on the files inside the shared folder.

https://github.com/owncloud/ocis/pull/8529  
https://github.com/cs3org/reva/pull/4534  
https://github.com/owncloud/ocis/issues/8273
