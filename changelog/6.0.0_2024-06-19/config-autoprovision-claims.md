Enhancement: Configurable claims for auto-provisioning user accounts

We introduce the new environment variables
"PROXY_AUTOPROVISION_CLAIM_USERNAME", "PROXY_AUTOPROVISION_CLAIM_EMAIL", and
"PROXY_AUTOPROVISION_CLAIM_DISPLAYNAME" which can be used to configure the
OIDC claims that should be used for auto-provisioning user accounts.

The automatic fallback to use the 'email' claim value as the username when
the 'preferred_username' claim is not set, has been removed.

Also it is now possible to autoprovision users without an email address.

https://github.com/owncloud/ocis/pull/8952
https://github.com/owncloud/ocis/issues/8635
https://github.com/owncloud/ocis/issues/6909
