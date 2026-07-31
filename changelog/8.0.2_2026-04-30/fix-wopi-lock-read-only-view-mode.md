Bugfix: Return 200 OK for WOPI Lock requests in read-only and view-only modes

OnlyOffice sends a WOPI Lock request when opening any document, even when
the user only has read access. The WOPI Lock handler was attempting to acquire
a CS3 write lock regardless of the view mode, causing a permission error for
read-only tokens that OnlyOffice displayed as an error message on load.

The Lock handler now returns 200 OK immediately for READ_ONLY and VIEW_ONLY
view modes without attempting to acquire a lock, consistent with the WOPI spec.

https://github.com/owncloud/ocis/pull/12257
