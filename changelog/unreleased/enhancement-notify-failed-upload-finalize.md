Enhancement: Notify the user when an upload fails to finalize

When an async blob commit fails after postprocessing and the node is reverted, the uploading user now receives a persisted "Upload failed" notification identifying the file, and the stale upload row is removed from the file list instead of the failure being silent.

https://github.com/owncloud/ocis/pull/
