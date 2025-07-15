Enhancement: Add Sharing NG endpoints

We've added new sharing ng endpoints to the graph beta api.
The following endpoints are added:

* /v1beta1/me/drive/sharedByMe
* /v1beta1/me/drive/sharedWithMe
* /v1beta1/roleManagement/permissions/roleDefinitions
* /v1beta1/roleManagement/permissions/roleDefinitions/{roleID}
* /v1beta1/drives/{drive-id}/items/{item-id}/createLink (create a sharing link)

https://github.com/owncloud/ocis/pull/7633
https://github.com/owncloud/ocis/pull/7686
https://github.com/owncloud/ocis/pull/7684
https://github.com/owncloud/ocis/pull/7683
https://github.com/owncloud/ocis/pull/7239
https://github.com/owncloud/ocis/pull/7687
https://github.com/owncloud/ocis/pull/7751
https://github.com/owncloud/libre-graph-api/pull/112
https://github.com/owncloud/ocis/issues/7436
https://github.com/owncloud/ocis/issues/6993
