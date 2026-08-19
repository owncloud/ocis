Enhancement: Add ResourceID field to UploadReady event

The UploadReady NATS event now includes a `ResourceID` field containing the
file's actual resource identifier (with the correct node OpaqueId). Previously,
only `FileRef` was available, whose `ResourceId.OpaqueId` is set to the space
root ID (required for CS3 gateway path resolution). Consumers that need the
file's unique identifier for Graph API or WebDAV operations can now use
`ResourceID.OpaqueId` directly.

https://github.com/owncloud/ocis/pull/12060
https://github.com/owncloud/ocis/issues/12056
https://github.com/owncloud/reva/pull/547
https://github.com/owncloud/reva/pull/560
