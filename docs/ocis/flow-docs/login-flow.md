---
title: "Login Flow"
date: 2020-05-04T20:47:00+01:00
weight: 43
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/flow-docs
geekdocFilePath: login-flow.md
---


## Login Flow

The following sequence diagram describes the [openid connect auth code flow](https://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth). The eight numbered steps and notes correspond to the [openid connect auth code flow steps](https://openid.net/specs/openid-connect-core-1_0.html#CodeFlowSteps). Example requests are based on the spec as well:

{{< mermaid class="text-center">}}
sequenceDiagram
    %% we have comments!! \o/
    %% this documents the login workflow
    %% examples taken from the oidc spec https://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth
    %% TODO add PKCE, see https://developer.okta.com/blog/2019/08/22/okta-authjs-pkce#use-pkce-to-make-your-apps-more-secure
    participant user as User
    participant client as Client
    participant proxy as ocis-proxy
    participant idp as IdP
    participant idm as LibreIDM
    participant ldap as External User Directory

    user->>+client: What is the content of my home?
        client->>+proxy: PROPFIND <br> no (or expired) auth
        Note over client,proxy: ocis needs to know the IdP that is<br>used to authenticate users. The<br>proxy will redirect unauthenticated<br>requests to that IdP.
        proxy-->>-client: 401 Unauthorized
        client->>+proxy: 1. The client starts a new openIDConnect Flow
        Note over client, proxy: GET /.well-known/openid-configuration
        proxy-->>-client: Return openidConnect configuration for the IdP
        client-->>client: 2. Client prepares an Authentication Request<br>containing the desired request parameters <br> and generates the code challenge (PKCE).
        client->>+idp: 3. Client sends the request and the code challenge to the Authorization Server.
        Note over client, idp: GET /authorize?<br>flow=oidc&response_type=code<br>&scope=openid%20profile%20email<br>&code_challenge=Y2SGoq9vtAp7YAavTaO0B550H_Rsj9DypiL7xZuFjOE<br>&code_challenge_method=S25&client_id=s6BhdRkqt3<br>&state=af0ifjsldkj<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb HTTP/1.1<br>Host: server.example.com
        Note over user, idp: 3. Authorization Server Authenticates the End-User.
        alt all users managed by idp/ocis idm
            idp->>+idm: LDAP query/bind
            idm-->>-idp: LDAP result
            Note over idp,ldap: In case  users are managed<br>in an external ldap they have to be<br> autoprovisioned in the ocis IdM<br> when they are loggin in.
        else all users authenticated by an external idp
            idp->>+ldap: Lookup of the user in the directory
            ldap-->>-idp: Lookup result
        end
        idp-->>-user: Idp presents the user an authentication prompt.
        user->>+idp: 5. User authenticates and gives consent.
        idp-->>-client: 6. Authorization Server sends the End-User back<br>to the Client with an Authorization Code.
        Note over client, idp: HTTP/1.1 302 Found<br>Location: https://client.example.org/cb?<br>code=SplxlOBeZQQYbYS6WxSbIA&state=af0ifjsldkj
        client->>+idp: 7. Client requests a response using the<br>Authorization Code and the code verifier at the Token Endpoint.
        Note over client, idp: POST /token HTTP/1.1<br>Host: server.example.com<br>Content-Type: application/x-www-form-urlencoded<br>grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA<br>&redirect_uri=https%3A%2F%2Fclient.example.org<br>&code_verifier=a98ccbe253754259963e6e2b67b5a044929446d7a15046cc8e3194022ad061d9d667dce91876418d9e6fe9f54819332e
        idp->>+idp: 8. IdP checks the code verifier (PKCE)
        idp-->>-client: 9. Client receives a response that contains an<br>ID Token and Access Token in the response body. <br>If offline access is requested, the client also receives a refresh token.
        Note over client, idp:  HTTP/1.1 200 OK<br>Content-Type: application/json<br>Cache-Control: no-store<br>Pragma: no-cache<br>{<br>"access_token": "SlAV32hkKG",<br>"token_type": "Bearer",<br>"refresh_token": "8xLOxBtZp8",<br>"expires_in": 3600,<br>"id_token": "a ... b.c ... d.e ... f" // must be a JWT<br>}
        client-->>client: 10. Client validates the ID token and<br>retrieves the End-User's Subject Identifier.
        client->>+proxy: PROPFIND <br> With access token
        proxy-->>-client: 207 Multi-Status
    client-->>-user: List of Files X, Y, Z ...
{{< /mermaid >}}
