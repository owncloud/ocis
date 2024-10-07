Bugfix: Avoid re-creating thumbnails

We fixed a bug that caused the system to re-create thumbnails for images, even
if a thumbnail already existed in the cache.

https://github.com/owncloud/ocis/pull/10251
