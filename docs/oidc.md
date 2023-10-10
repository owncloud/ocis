# Ocis as a OIDC Identity Provider
Ocis can function as an Identity Provider (IdP) for OpenID Connect (OIDC) authentication .
OIDC defines a discovery mechanism, called OpenID Connect Discovery,
where an OpenID server publishes its metadata at a well-known URL, typically

`https://server.com/.well-known/openid-configuration`

This URL returns a JSON listing of the OpenID/OAuth endpoints, supported scopes and claims, public keys used to sign the tokens, and other details.
The clients can use this information to construct a request to the OpenID server.
The field names and values are defined in the [OpenID Connect Discovery Specification](https://openid.net/specs/openid-connect-discovery-1_0.html).
Here is an example of data returned:
```json
{
  "issuer": "https://localhost:9200",
  "authorization_endpoint": "https://localhost:9200/signin/v1/identifier/_/authorize",
  "token_endpoint": "https://localhost:9200/konnect/v1/token",
  "userinfo_endpoint": "https://localhost:9200/konnect/v1/userinfo",
  "end_session_endpoint": "https://localhost:9200/signin/v1/identifier/_/endsession",
  "check_session_iframe": "https://localhost:9200/konnect/v1/session/check-session.html",
  "jwks_uri": "https://localhost:9200/konnect/v1/jwks.json",
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

## Authorization code grant
1. Requesting authorization
To initiate the OIDC Code Flow, you can use tools like curl and a web browser.
The user should be directed to a URL similar to the following to authenticate
and give their consent (bypassing consent is against the standard):

```plaintext
https://localhost:9200/signin/v1/identifier/_/authorize?client_id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69&scope=openid+profile+email+offline_access&response_type=code&redirect_uri=http%3A%2F%2Flocalhost%2F
```

After a successful authentication, the browser will redirect to a URL that looks like this:

```plaintext
http://localhost/?code=mfWsjEL0mc8gx0ftF9LFkGb__uFykaBw&scope=openid%20profile%20email%20offline_access&session_state=32b08dd722f1227f0bd635d98e36d5d066eba34c2d1cbed6447c5e7085608ceb.Q8YK6cXpZVDgiB3nFiVR2OOtJffu5AzJRXTMSDX3KXM&state=
```

For the next step extract the code from the URL.
In the above example,
the code is `mfWsjEL0mc8gx0ftF9LFkGb__uFykaBw`

2. Requesting an access token
The next step in the  OIDC Code Flow involves HTTP POST request
to the token endpoint of the **Ocis Identity Server**.

```bash
curl -vk -X POST https://localhost:9200/konnect/v1/token
-H "Content-Type: application/x-www-form-urlencoded"
-d "grant_type=authorization_code
&code=3a3PTcO-WWXfN3l1mDN4u7G5PzWFxatU
&redirect_uri=http%3A%2F%2Flocalhost%2Fowncloud%2FDevelopment%2Fpraticeground%2Fsrc%2Fdemo.php
&client_id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69
&client_secret=UBntmLjC2yYCeHwsyj73Uwo9TAaecAetRwMw0xYcvNL9yRdLSUi0hUAHfvCHFeFh"
```

 - Code exchange response
```json
{
  "access_token": "eyJhbGciOid...",
  "token_type": "Bearer",
  "id_token": "eyJhbGciOi...",
  "refresh_token": "eyJhbGciOiJ...",
  "expires_in": 300
}
```

3. Refreshing an access token`
If the access token has expired, you can get a new one with the refresh token.
```bash
curl -vk -X POST https://localhost:9200/konnect/v1/token
  -d "grant_type=refresh_token" \
  -d "refresh_token=eyJhbGciOiJ..." \
  -d "redirect_uri=redirect_url_path" \
  -d "client_id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69" \
  -d "client_secret=UBntmLjC2yYCeHwsyj73Uwo9TAaecAetRwMw0xYcvNL9yRdLSUi0hUAHfvCHFeFh" \
```
- Refreshing an access token response
```json
{
  "access_token": "eyJhbGciOi...",
  "token_type": "Bearer",
  "expires_in": 300
}
```
