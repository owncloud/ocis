Bugfix: Fixed `sharedWithMe` response for OCM shares

OCM shares returned in the `sharedWithMe` response did not have the `mimeType` property
populated correctly.

https://github.com/owncloud/ocis/pull/10501
https://github.com/owncloud/ocis/issues/10495
