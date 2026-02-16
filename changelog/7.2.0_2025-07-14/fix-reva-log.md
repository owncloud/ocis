Bugfix: Fix the reva log interceptor

Fix the reva log interceptor. Implemented the Unwrap interface to allow TUS middleware to handle correctly
SetReadDeadline and SetWriteDeadline functions and to avoid the error during the upload.

https://github.com/owncloud/ocis/pull/11348    
https://github.com/owncloud/ocis/issues/10857
