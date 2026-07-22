Enhancement: Add audience restriction for OIDC access tokens

We added a new proxy configuration option `PROXY_OIDC_ACCESS_TOKEN_VERIFY_AUD`
which lets operators restrict the JWT access tokens that are accepted to a list
of allowed audiences. When set, a token is only accepted if one of the
configured values is present in its `aud` (audience) claim or matches its `azp`
(authorized party) claim. This is useful when the OIDC provider (IDP) is shared
between multiple applications, so that tokens issued for another application are
not accepted by ownCloud Infinite Scale.

The check is disabled by default (empty list) to preserve the existing behavior.

https://github.com/owncloud/ocis/pull/12634

