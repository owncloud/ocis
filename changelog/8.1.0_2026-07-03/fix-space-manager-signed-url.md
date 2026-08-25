Bugfix: Fix space management middleware removing users from spaces on download

The space management middleware ran on every authenticated request, including
signed URL requests used for file downloads. Since signed URL auth does not
carry OIDC claims, the middleware interpreted the absence of claims as "user
should have no space access" and removed the user from all project spaces.
On the next OIDC request the user was re-added, causing an oscillating
add/remove cycle that led to intermittent download failures and transient
"space not found" errors.

The middleware now skips reconciliation entirely when no OIDC claims are
present in the request context.

https://github.com/owncloud/ocis/pull/12285
https://github.com/owncloud/ocis/issues/12285
