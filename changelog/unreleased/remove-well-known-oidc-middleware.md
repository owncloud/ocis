Enhancement: Introduce staticroutes package & remove well-known OIDC middleware

We have introduced a new static routes package to the proxy. This package
is responsible for serving static files and oidc well-known endpoint `/.well-known/openid-configuration`.
We have removed the well-known middleware for OIDC and moved it
to the newly introduced static routes module in the proxy.

https://github.com/owncloud/ocis/issues/6095
https://github.com/owncloud/ocis/pull/8541

