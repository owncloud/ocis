Bugfix: Return 404 when deleting a non-existing link share on a space

Deleting an already-removed link share via
`DELETE /graph/v1beta1/drives/{driveID}/root/permissions/{permissionID}`
returned `400 Bad Request` instead of `404 Not Found`. The underlying
reva `RemovePublicShare` and `RemoveShare` handlers now propagate
not-found errors as `CODE_NOT_FOUND` rather than `CODE_INTERNAL`,
ensuring the correct `404` HTTP response is returned.

https://github.com/owncloud/ocis/issues/12266
