Bugfix: CheckFileInfo will return a 404 error if the target file isn't found

Previously, the request failed with a 500 error code, but it it will fail with a 404 error code

https://github.com/owncloud/ocis/pull/10112
