Bugfix: Fix claim selector based routing for basic auth

We've fixed the claim selector based routing for requests using basic auth.
Previously requests using basic auth have always been routed to the DefaultPolicy when using the claim selector despite the set cookie because the basic auth middleware fakes some OIDC claims.

Now the cookie is checked before routing to the DefaultPolicy and therefore set cookie will also be respected for requests using basic auth.

https://github.com/owncloud/ocis/pull/2779
