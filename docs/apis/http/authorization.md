---
title: Authorization
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/http/
geekdocFilePath: authorization.md
---

In its default configuration, Infinite Scale supports three authentication methods as outlined on the [OIDC official site](https://openid.net/specs/openid-connect-core-1_0.html#rfc.section.3):
1. Authorization Code Flow
2. Implicit Flow
3. Hybrid Flow

For detailed information on Infinite Scale's support for OpenID Connect (OIDC), please consult the [OIDC section](https://owncloud.dev/ocis/identity-provider/oidc).
To authenticate a client app using OIDC, both `client_id` and `client_secret` are essential. Infinite Scale does not offer dynamic registration. The required data for the default [ownCloud clients](https://doc.owncloud.com/server/next/admin_manual/configuration/user/oidc/oidc.html#client-ids-secrets-and-redirect-uris) can be found in the link and are availble for the following apps:
- Desktop
- Android
- iOS

While selecting an ownCloud client for authentication, take note of specific limitations such as the `Redirect URI`:

| Source | Redirect URI |
|------|--------|
|Android|oc://android.owncloud.com|
|iOS|oc://ios.owncloud.com|
|Desktop|http://127.0.0.1 <br> http://localhost |

In this example, the desktop app's `client_id` and `client_secret` are being used.

```bash
client_id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69
client_secret=UBntmLjC2yYCeHwsyj73Uwo9TAaecAetRwMw0xYcvNL9yRdLSUi0hUAHfvCHFeFh
```

## Authorization Code Flow

1. Requesting authorization

   To initiate the OIDC Code Flow, you can use tools like curl and a web browser.
   The user should be directed to a URL to authenticate and give their consent (bypassing consent is against the standard):

    ```plaintext
    https://ocis.test/signin/v1/identifier/_/authorize?client_id=client_id&scope=openid+profile+email+offline_access&response_type=code&redirect_uri=http://path-to-redirect-uri
    ```

    After a successful authentication, the browser will redirect to a URL that looks like this:

    ```plaintext
    http://path-to-redirect-uri?code=mfWsjEL0mc8gx0ftF9LFkGb__uFykaBw&scope=openid%20profile%20email%20offline_access&session_state=32b08dd...&state=
    ```

    For the next step extract the code from the URL.

    In the above example,
    the code is `mfWsjEL0mc8gx0ftF9LFkGb__uFykaBw`

2. Requesting an access token

   The next step in the OIDC Code Flow involves an HTTP POST request
   to the token endpoint of the **Infinite Scale Identity Server**.

    ```bash
    curl -vk -X POST https://ocis.test/konnect/v1/token \
    -d "grant_type=authorization_code" \
    -d "code=3a3PTcO-WWXfN3l1mDN4u7G5PzWFxatU" \
    -d "redirect_uri=http:path-to-redirect-uri" \
    -d "client_id=client_id" \
    -d "client_secret=client_secret"
    ```

   Response looks like this:
    ```json
    {
    "access_token": "eyJhbGciOid...",
    "token_type": "Bearer",
    "id_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOiJ...",
    "expires_in": 300
    }
    ```

3. Refreshing an access token

   If the access token has expired, you can get a new one using the refresh token.
    ```bash
    curl -vk -X POST https://ocis.test/konnect/v1/token \
    -d "grant_type=refresh_token" \
    -d "refresh_token=eyJhbGciOiJ..." \
    -d "redirect_uri=http://path-to-redirect-uri" \
    -d "client_id=client_id" \
    -d "client_secret=client_secret"
    ```

   Response looks like this:
    ```json
    {
    "access_token": "eyJhbGciOi...",
    "token_type": "Bearer",
    "expires_in": 300
    }
    ```

## Implicit Code Flow

When using the implicit flow, tokens are provided in a URI fragment of the redirect URL.
Valid values for the `response_type` request parameter are:

- token
- id_token token

{{< hint type=warning title="Important Warning" >}}
If you are using the implicit flow, `nonce` parameter is required in the initial `/authorize` request.
`nonce=pL3UkpAQPZ8bTMGYOmxHY/dQABin8yrqipZ7iN0PY18=`

bash command to generate cryptographically random value
```bash
openssl rand -base64 32
```
{{< /hint >}}

The user should be directed to a URL to authenticate and give their consent (bypassing consent is against the standard):
```bash
https://ocis.test/signin/v1/identifier/_/authorize?client_id=client_id&scope=openid+profile+email+offline_access&response_type=id_token+token&redirect_uri=http://path-to-redirect-uri&nonce=pL3UkpAQPZ8bTMGYOmxHY/dQABin8yrqipZ7iN0PY18=
 ```

After a successful authentication, the browser will redirect to a URL that looks like this:
```bash
http://path-to-redirect-uri#access_token=eyJhbGciOiJQUzI...&expires_in=300&id_token=eyJhbGciOiJ...&scope=email%20openid%20profile&session_state=c8a1019f5e054d...&state=&token_type=Bearer
```

For the next step, extract the access_token from the URL.
```bash
access_token = 'eyJhbGciOiJQ...'
 ```

## Hybrid Flow
The Hybrid Flow in OpenID Connect melds features from both the Implicit and Authorization Code flows. It allows clients to directly retrieve certain tokens from the Authorization Endpoint, yet also offers the option to acquire additional tokens from the Token Endpoint.

The Authorization Server redirects back to the client with appropriate parameters in the response, based on the value of the response_type request parameter:
- code token
- code id_token
- code id_token token
