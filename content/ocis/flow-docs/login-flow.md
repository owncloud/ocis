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
    participant glauth as ocis-glauth
    participant graph as ocis-graph
    participant accounts as ocis-accounts
    participant ldap as external LDAP server

    user->>+client: What is the content of my home?

        client->>+proxy: PROPFIND <br> no (or expired) auth
        Note over client,proxy: ocis needs to know the IdP that is<br>used to authenticate users. The<br>proxy will redirect unauthenticated<br>requests to that IdP.
        proxy-->>-client: 302 Found
        Note over client, idp: HTTP/1.1 302 Found<br>Location: https://server.example.com/authorize?<br>response_type=code&<br>scope=openid%20profile%20email<br>&client_id=s6BhdRkqt3<br>&state=af0ifjsldkj<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb

        Note over client, idp: We should follow the OpenID Connect Discovery protocol
        Note over client, idp: Clients might fall back to the ocis server if the discovery failed.<br>We can provide a webfinger endpoint there to let guests use an idp<br>that is backed by the accounts service.
        Note over client, idp: For now, clients can only handle one IdP, which is configured in ocis.

        client-->>client: 1. Client prepares an Authentication Request<br>containing the desired request parameters.

        client->>+idp: 2. Client sends the request to the Authorization Server.
        Note over client, idp: GET /authorize?<br>response_type=code<br>&scope=openid%20profile%20email<br>&client_id=s6BhdRkqt3<br>&state=af0ifjsldkj<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb HTTP/1.1<br>Host: server.example.com
        Note over user, idp: 3. Authorization Server Authenticates the End-User.
        Note over idp,ldap: Either an IdP already exists or a new one is introduced. Since we are not yet using oidc discovery we can only use one IdP.
        alt all users managed by idp/ocis
            idp->>+glauth: LDAP query/bind
            glauth->>+graph: GET user with Basic Auth<br>GraphAPI
            graph->>+accounts: internal GRPC
            accounts-->>-graph: response
            graph-->>-glauth: OData response
            glauth-->>-idp: LDAP result
            Note over accounts,ldap: In case internal users are managed<br>in an external ldap they have to be<br>synced to the accounts service to<br>show up as recipients during sharing.
        else all users authenticated by an external idp
            idp->>+ldap: LDAP query/bind
            ldap-->>-idp: LDAP result
            alt guest accounts managed in ocis / lookup using glauth proxy:
                Note over idp,glauth: Idp is configured to use glauth as a<br>second ldap server.
                idp->>+glauth: LDAP query/bind
                glauth->>+graph: GET user with Basic Auth<br>GraphAPI
                graph->>+accounts: internal GRPC
                accounts-->>-graph: response
                graph-->>-glauth: OData response
                glauth-->>-idp: LDAP result
            else guest account provisioned by other means
                Note over accounts, ldap: In case guest accounts are managed<br>in an existing ldap they need to be<br>synced to the accounts service to<br>be able to login and show up as<br>recipients during sharing.
            end
        end
        Note over user, idp: 4. Authorization Server obtains End-User Consent/Authorization.
        idp-->>-client: 5. Authorization Server sends the End-User back<br>to the Client with an Authorization Code.
        Note over client, idp: HTTP/1.1 302 Found<br>Location: https://client.example.org/cb?<br>code=SplxlOBeZQQYbYS6WxSbIA&state=af0ifjsldkj

        client->>+idp: 6. Client requests a response using the<br>Authorization Code at the Token Endpoint.
        Note over client, idp: POST /token HTTP/1.1<br>Host: server.example.com<br>Content-Type: application/x-www-form-urlencoded<br>grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb
        idp-->>-client: 7. Client receives a response that contains an<br>ID Token and Access Token in the response body.
        Note over client, idp:  HTTP/1.1 200 OK<br>Content-Type: application/json<br>Cache-Control: no-store<br>Pragma: no-cache<br>{<br>"access_token": "SlAV32hkKG",<br>"token_type": "Bearer",<br>"refresh_token": "8xLOxBtZp8",<br>"expires_in": 3600,<br>"id_token": "a ... b.c ... d.e ... f" // must be a JWT<br>}


        client-->>client: 8. Client validates the ID token and<br>retrieves the End-User's Subject Identifier.

        client->>+proxy: PROPFIND <br> With access token
        proxy-->>-client: 207 Multi-Status
    client-->>-user: List of Files X, Y, Z ...
{{< /mermaid >}}
