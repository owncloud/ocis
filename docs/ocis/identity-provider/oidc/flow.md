---
title: Flow
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/identity-provider/oidc
geekdocFilePath: flow.md
---

In Infinite Scale, authentication can follow one of the three methods described on the [official site](https://openid.net/specs/openid-connect-core-1_0.html#rfc.section.3):
1. Authorization Code Flow
2. Implicit Flow
3. Hybrid Flow

To authenticate using OIDC, both `client_id` and `client_secret` are essential. For oidc request, desktop-client `client_id` and `client_secret` can be used.
```bash
client_id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69
client_secret=UBntmLjC2yYCeHwsyj73Uwo9TAaecAetRwMw0xYcvNL9yRdLSUi0hUAHfvCHFeFh
```
For more specifics, refer to the [ownCloud documentation](https://doc.owncloud.com/server/next/admin_manual/configuration/user/oidc/oidc.html#client-ids-secrets-and-redirect-uris)

# Authentication Code Flow
1. Requesting authorization

   To initiate the OIDC Code Flow, you can use tools like curl and a web browser.
   The user should be directed to a URL to authenticate and give their consent (bypassing consent is against the standard):

    ```plaintext
    https://ocis.test:9200/signin/v1/identifier/_/authorize?client_id=client_id&scope=openid+profile+email+offline_access&response_type=code&redirect_uri=http://path-to-redirect-uri
    ```

    After a successful authentication, the browser will redirect to a URL that looks like this:

    ```plaintext
    http://path-to-redirect-uri?code=mfWsjEL0mc8gx0ftF9LFkGb__uFykaBw&scope=openid%20profile%20email%20offline_access&session_state=32b08dd...&state=
    ```

    For the next step extract the code from the URL.

    In the above example,
    the code is `mfWsjEL0mc8gx0ftF9LFkGb__uFykaBw`

2. Requesting an access token

   The next step in the  OIDC Code Flow involves an HTTP POST request
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

   If the access token has expired, you can get a new one with the refresh token.
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

# Implicit Code Flow
   In implicit flow, tokens return via the URI fragment that has been viewed as less secure than other flows.
   Value of the `response_type` request parameter could be :
   - token
   - id_token token

   > **Note**
   >
   > If you are using the implicit flow, `nonce` parameter is required in the initial `/authorize` request,
   > nonce=8e641aff9b22e3f0c6d052b6b443a3ac

   ```bash
    https://ocis.test/signin/v1/identifier/_/authorize?client_id=client_id&scope=openid+profile+email+offline_access&response_type=id_token+token&redirect_uri=http://path-to-redirect-uri&nonce=8e641aff9b22e3f0c6d052b6b443a3ac
   ```

   After a successful authentication, the browser will redirect to a URL that looks like this:
   ```bash
    http://path-to-redirect-uri#access_token=eyJhbGciOiJQUzI...&expires_in=300&id_token=eyJhbGciOiJ...&scope=email%20openid%20profile&session_state=c8a1019f5e054d...&state=&token_type=Bearer
   ```

   For the next step extract the access_token from the URL.
   ```bash
   access_token = 'eyJhbGciOiJQ...'
   ```
# Hybrid Flow
   The Hybrid Flow in OpenID Connect melds features from both the Implicit and Authorization Code flows. It allows clients to directly retrieve certain tokens from the Authorization Endpoint, yet also offers the option to acquire additional tokens from the Token Endpoint.

   The Authorization Server redirects back to the client with appropriate parameters in the response, based on the value of the response_type request parameter:
   - code token
   - code id_token
   - code id_token token
