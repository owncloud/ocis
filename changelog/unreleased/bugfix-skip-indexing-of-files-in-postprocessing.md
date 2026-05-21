Bugfix: Skip indexing of files still in postprocessing

When the search service re-indexed a space in response to an UploadReady
event, the walker visited sibling nodes whose blobs were not yet finalized
in the blobstore. Content extraction for those in-flight nodes triggered
spurious storage-users error logs (S3 NoSuchKey). The walker now skips
nodes marked as processing; they are indexed when their own UploadReady
event arrives.

https://github.com/owncloud/ocis/pull/00000
