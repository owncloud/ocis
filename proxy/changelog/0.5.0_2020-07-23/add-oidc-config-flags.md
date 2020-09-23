Change: Add OIDC config flags

To authenticate requests with an oidc provider we added two environment variables:
- `PROXY_OIDC_ISSUER="https://localhost:9200"` and
- `PROXY_OIDC_INSECURE=true`

This changes ocis-proxy to now load the oidc-middleware by default, requiring a bearer token and exchanging the email in the OIDC claims for an account id at the ocis-accounts service.

Setting `PROXY_OIDC_ISSUER=""` will disable the OIDC middleware.

https://github.com/owncloud/ocis/proxy/pull/66
