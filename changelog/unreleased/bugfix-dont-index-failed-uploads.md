Bugfix: Don't index failed uploads

The search service was indexing uploads even when they failed. This caused
unnecessary index operations for incomplete or errored file transfers. The fix
skips indexing when the UploadReady event indicates the upload has failed.

https://github.com/owncloud/ocis/pull/12121
