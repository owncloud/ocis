Bugfix: graph/sharedWithMe align IDs with webdav response

The IDs of the driveItems returned by the 'graph/v1beta1/me/drive/sharedWithMe'
endpoint are now aligned with the IDs returned in the PROPFIND response
of the webdav service.

https://github.com/owncloud/ocis/pull/8467
https://github.com/owncloud/ocis/issues/8420
https://github.com/owncloud/ocis/issues/8080
