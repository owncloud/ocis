Enhancement: Introduce OpenID Connect middleware

Added an openid connect middleware that will try to authenticate users using
OpenID Connect. The claims will be added to the context of the request.

https://github.com/owncloud/ocis-pkg/issues/8
