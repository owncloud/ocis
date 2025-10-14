Enhancement: Support spec violating AD FS access token issuer

AD FS `/adfs/.well-known/openid-configuration` has an optional `access_token_issuer` which, in violation of the OpenID Connect spec, takes precedence over `issuer`.

https://github.com/owncloud/ocis/pull/7138
