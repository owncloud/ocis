Enhancement: Notify the user when an upload fails to finalize

When an async blob commit fails after postprocessing and the node is reverted, the uploading user now receives a persisted "Upload failed" notification identifying the file, and the stale upload row is removed from the file list instead of the failure being silent. This notification respects the same in-app notification setting as the antivirus and policy postprocessing notifications, and is suppressed when one of those more specific notifications already covered the same upload, so a rejected file no longer produces two separate notifications for the same failure.

https://github.com/owncloud/ocis/pull/12987
