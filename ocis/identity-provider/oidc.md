---
title: "OIDC"
date: 2023-10-10T00:00:00+00:00
weight: 21
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/identity-provider
geekdocFilePath: oidc.md
geekdocCollapseSection: true
---

Infinite Scale has implemented OpenID Connect (OIDC) for authentication.
OIDC defines a discovery mechanism, called OpenID Connect Discovery,
where an OpenID server publishes its metadata at a well-known URL, typically:

`https://ocis.test/.well-known/openid-configuration`

This URL returns a JSON listing of the OpenID/OAuth endpoints, supported scopes and claims, public keys used to sign the tokens, and other details.
The clients can use this information to construct a request to the OpenID server.
The field names and values are defined in the [OpenID Connect Discovery Specification](https://openid.net/specs/openid-connect-discovery-1_0.html).
Here is an example of data returned:
```json
{
  "issuer": "https://ocis.test",
  "authorization_endpoint": "https://ocis.test/signin/v1/identifier/_/authorize",
  "token_endpoint": "https://ocis.test/konnect/v1/token",
  "userinfo_endpoint": "https://ocis.test/konnect/v1/userinfo",
  "end_session_endpoint": "https://ocis.test/signin/v1/identifier/_/endsession",
  "check_session_iframe": "https://ocis.test/konnect/v1/session/check-session.html",
  "jwks_uri": "https://ocis.test/konnect/v1/jwks.json",
  "scopes_supported": [
    "openid",
    "offline_access",
    "profile",
    "email",
    "LibgreGraph.UUID",
    "LibreGraph.RawSub"
  ],
  "response_types_supported": [
    "id_token token",
    "id_token",
    "code id_token",
    "code id_token token"
  ],
  "subject_types_supported": [
    "public"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256",
    "RS384",
    "RS512",
    "PS256",
    "PS384",
    "PS512"
  ],
  "userinfo_signing_alg_values_supported": [
    "RS256",
    "RS384",
    "RS512",
    "PS256",
    "PS384",
    "PS512"
  ],
  "request_object_signing_alg_values_supported": [
    "ES256",
    "ES384",
    "ES512",
    "RS256",
    "RS384",
    "RS512",
    "PS256",
    "PS384",
    "PS512",
    "none",
    "EdDSA"
  ],
  "token_endpoint_auth_methods_supported": [
    "client_secret_basic",
    "none"
  ],
  "token_endpoint_auth_signing_alg_values_supported": [
    "RS256",
    "RS384",
    "RS512",
    "PS256",
    "PS384",
    "PS512"
  ],
  "claims_parameter_supported": true,
  "claims_supported": [
    "iss",
    "sub",
    "aud",
    "exp",
    "iat",
    "name",
    "family_name",
    "given_name",
    "email",
    "email_verified"
  ],
  "request_parameter_supported": true,
  "request_uri_parameter_supported": false
}
```

Refer to the [Authorization](https://owncloud.dev/apis/http/authorization) section for example usages and details.
