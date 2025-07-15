Enhancement: Add GetPermission API

Graph service: added GET /v1beta1/drives/{driveId}/items/{itemId}/permissions/{permissionId} (and space-root equivalent) so clients can retrieve a single permission instead of listing all.

https://github.com/owncloud/ocis/issues/8616
https://github.com/owncloud/ocis/pull/11477
