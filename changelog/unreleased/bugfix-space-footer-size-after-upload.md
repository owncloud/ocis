Bugfix: Update the total size below the file list of a Space

In the root of a project Space, the footer below the file list ("N items with X
in total") showed the size of the Space from the time the Spaces were loaded,
so it stayed stale after uploading, deleting or restoring files until the page
was reloaded. Web already fetches the Space again after these actions to update
its quota. It now updates the Space size (the used quota) from the same
response as well.

https://github.com/owncloud/ocis/issues/13259
