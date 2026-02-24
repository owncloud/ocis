Bugfix: Fix UploadReady event reporting wrong file resource ID

The UploadReady event and upload-finished callback both set
`FileRef.ResourceId.OpaqueId` to the space ID instead of the file's
actual node ID. Every upload event reported the space root UUID as
the file identifier. Current consumers are unaffected because they
resolve files by path, but any consumer relying on the resource ID
to identify the uploaded file would get the space root instead.

https://github.com/owncloud/ocis/issues/12056
https://github.com/owncloud/reva/pull/546
