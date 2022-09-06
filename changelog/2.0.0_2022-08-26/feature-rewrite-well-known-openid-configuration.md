Enhancement: Fix behavior for foobar (in present tense)

We've added the configuration option `PROXY_OIDC_REWRITE_WELLKNOWN` to rewrite the `/.well-known/openid-configuration` endpoint.
If active, it serves the `/.well-known/openid-configuration` response of the original IDP configured in `OCIS_OIDC_ISSUER` / `PROXY_OIDC_ISSUER`. This is needed so that the Desktop Client, Android Client and iOS Client can discover the OIDC identity provider.

Previously this rewrite needed to be performed with an external proxy as NGINX or Traefik if an external IDP was used.

https://github.com/owncloud/ocis/pull/4346
https://github.com/owncloud/ocis/issues/2819
https://github.com/owncloud/ocis/issues/3280
