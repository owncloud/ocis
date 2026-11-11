Bugfix: Fix AuthApp and basic auth requests with custom PROXY_USER_OIDC_CLAIM

We fixed AuthApp (app password) and HTTP basic auth requests returning an
internal server error whenever PROXY_USER_OIDC_CLAIM was set to anything
other than its default value. Both authenticators put the already
authenticated user directly into the request context instead of faking
incomplete claims for the account resolver to re-resolve. We also made sure
a user who only ever authenticates this way still gets a default role
assigned, and that these requests keep getting rejected on deployments with
multiple instances enabled, since neither carries tenant-membership
information.

https://github.com/owncloud/ocis/pull/13107
