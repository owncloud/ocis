Bugfix: Lookup user by id for presigned URLs

Phoenix will send the `userid`, not the `username` as the `OC-Credential` for presigned URLs. This PR uses the new `ocisid` claim in the OIDC userinfo to pass the userid to the account middleware.

https://github.com/owncloud/ocis-proxy/pull/85
https://github.com/owncloud/ocis-pkg/pull/50
https://github.com/owncloud/ocis/issues/436
