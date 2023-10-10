---
title: Authorization code flow
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/http/oidc
geekdocFilePath: authorization-code-flow.md
---

1. Requesting authorization

   To initiate the OIDC Code Flow, you can use tools like curl and a web browser.
   The user should be directed to a URL similar to the following to authenticate
   and give their consent (bypassing consent is against the standard):

    ```plaintext
    https://ocis.test:9200/signin/v1/identifier/_/authorize?client_id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69&scope=openid+profile+email+offline_access&response_type=code&redirect_uri=http%3A%2F%2Flocalhost%2F
    ```

    After a successful authentication, the browser will redirect to a URL that looks like this:

    ```plaintext
    http://ocis.test?code=mfWsjEL0mc8gx0ftF9LFkGb__uFykaBw&scope=openid%20profile%20email%20offline_access&session_state=32b08dd722f1227f0bd635d98e36d5d066eba34c2d1cbed6447c5e7085608ceb.Q8YK6cXpZVDgiB3nFiVR2OOtJffu5AzJRXTMSDX3KXM&state=
    ```

    For the next step extract the code from the URL.

    In the above example,
    the code is `mfWsjEL0mc8gx0ftF9LFkGb__uFykaBw`

2. Requesting an access token

   The next step in the  OIDC Code Flow involves HTTP POST request
   to the token endpoint of the **Ocis Identity Server**.

    ```bash
    curl -vk -X POST https://ocis.test:9200/konnect/v1/token \
    -d "grant_type=authorization_code" \
    -d "code=3a3PTcO-WWXfN3l1mDN4u7G5PzWFxatU" \
    -d "redirect_uri=http-path-to-redirect-uri" \
    -d "client_id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69" \
    -d "client_secret=UBntmLjC2yYCeHwsyj73Uwo9TAaecAetRwMw0xYcvNL9yRdLSUi0hUAHfvCHFeFh"
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
    curl -vk -X POST https://localhost:9200/konnect/v1/token \
    -d "grant_type=refresh_token" \
    -d "refresh_token=eyJhbGciOiJ..." \
    -d "redirect_uri=redirect_url_path" \
    -d "client_id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69" \
    -d "client_secret=UBntmLjC2yYCeHwsyj73Uwo9TAaecAetRwMw0xYcvNL9yRdLSUi0hUAHfvCHFeFh"
    ```

   - Refreshing an access token response
    ```json
    {
    "access_token": "eyJhbGciOi...",
    "token_type": "Bearer",
    "expires_in": 300
    }
    ```
