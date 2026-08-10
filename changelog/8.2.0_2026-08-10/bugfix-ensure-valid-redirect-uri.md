Bugfix: ensure the redirect URI for the IDP is valid

The URI sent as redirect URI for the IDP will be validated in oCIS.
Invalid URIs will return a 500 error. Note that this should never happen
under normal circumstances.

https://github.com/owncloud/ocis/pull/12444
https://github.com/owncloud/ocis/pull/12479
