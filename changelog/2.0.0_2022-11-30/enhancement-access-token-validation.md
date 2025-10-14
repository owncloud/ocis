Enhancement: Improve validation of OIDC access tokens

Previously OIDC access tokes were only validated by requesting the userinfo from
the IDP. It is now possible to enable additional verification if the IDP issues
access tokens in JWT format. In that case the oCIS proxy service will now verify
the signature of the token using the public keys provided by jwks_uri endpoint
of the IDP. It will also verify if the issuer claim (iss) matches the expected
values.

The new validation is enabled by setting `PROXY_OIDC_ACCESS_TOKEN_VERIFY_METHOD`
to "jwt". Which is also the default. Setting it to "none" will disable the feature.

https://github.com/owncloud/ocis/issues/3841
https://github.com/owncloud/ocis/pull/4227
