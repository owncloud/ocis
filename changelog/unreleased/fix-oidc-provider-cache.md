Bugfix: Fix the OIDC provider cache

We've fixed the OIDC provider cache. It never had a cache hit before this fix.
Under some circumstances it could cause a painfully slow OCIS if the IDP wellknown endpoint takes some time to respond.

https://github.com/owncloud/ocis/pull/4600
