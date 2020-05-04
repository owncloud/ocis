---
title: "Login Flow"
date: 2020-05-04T20:47:00+01:00
weight: 43
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: login-flow.md
---


## Login Flow

The following sequence diagram describes the [openid connect auth flow](https://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth):

{{< mermaid class="text-center">}}
sequenceDiagram
    %% we have comments!! \o/
    %% this documents the login workflow
    %% examples taken from the oidc spec https://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth
    %% TODO add PKCE, see https://developer.okta.com/blog/2019/08/22/okta-authjs-pkce#use-pkce-to-make-your-apps-more-secure
    participant user as User
    participant client as Client
    participant proxy as ocis-proxy
    participant idp as external IdP
    participant konnectd as ocis-konnectd IdP
    participant glauth as ocis-glauth
    participant graph as ocis-graph
    participant accounts as ocis-accounts
    participant ldap as external LDAP server

    user->>+client: What is the content of my home?

        client->>+proxy: PROPFIND <br> no (or expired) auth
        alt proxy can decide which idp to use
          Note over client, idp: We may not be able to differentiate guests from users
          proxy-->>client: 302 Found
          Note over client, idp: HTTP/1.1 302 Found<br>Location: https://server.example.com/authorize?<br>response_type=code&<br>scope=openid%20profile%20email<br>&client_id=s6BhdRkqt3<br>&state=af0ifjsldkj<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb
        else client needs to discover the idp
          proxy-->>-client: 401 Unauthorized
          Note over client, idp: Follow OpenID Connect Discovery protocol
          Note over client, idp: Clients might fall back to the ocis server if the discovery failed.<br>We can provide a webfinger endpoint there to let guests use an idp<br>that is backed bythe accounts service.
          Note over client, idp: For now, always use ocis well known endpoint to discover idp?<br>We can check the email in the accounts service.
        end

        client-->>client: 1. Client prepares an Authentication Request<br>containing the desired request parameters.

        alt all users authenticated by an external idp
            client->>+idp: 2. Client sends the request to the Authorization Server.
            Note over client, idp: GET /authorize?<br>response_type=code<br>&scope=openid%20profile%20email<br>&client_id=s6BhdRkqt3<br>&state=af0ifjsldkj<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb HTTP/1.1<br>Host: server.example.com
            Note over user, idp: 3. Authorization Server Authenticates the End-User.
            idp->>+ldap: LDAP query/bind
            ldap-->>-idp: LDAP result
            alt guest accounts managed in ocis / lookup using glauth proxy:
                idp->>+glauth: LDAP query/bind
                glauth->>+graph: GET user with Basic Auth<br>GraphAPI
                graph->>+accounts: internal GRPC
                accounts-->>-graph: response
                graph-->>-glauth: OData response
                glauth-->>-idp: LDAP result
            else guest account provisioned by other means
                Note over idp, ldap: In case guest accounts are stored in an existing ldap they need to be synced to the accounts service to show up as recipients during sharing.
            end

            Note over user, idp: 4. Authorization Server obtains End-User Consent/Authorization.
            idp-->>-client: 5. Authorization Server sends the End-User back<br>to the Client with an Authorization Code.
            Note over client, idp: HTTP/1.1 302 Found<br>Location: https://client.example.org/cb?<br>code=SplxlOBeZQQYbYS6WxSbIA&state=af0ifjsldkj

            client->>+idp: 6. Client requests a response using the<br>Authorization Code at the Token Endpoint.
            Note over client, idp: POST /token HTTP/1.1<br>Host: server.example.com<br>Content-Type: application/x-www-form-urlencoded<br>grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb
            idp-->>-client: 7. Client receives a response that contains an<br>ID Token and Access Token in the response body.
            Note over client, idp:  HTTP/1.1 200 OK<br>Content-Type: application/json<br>Cache-Control: no-store<br>Pragma: no-cache<br>{<br>"access_token": "SlAV32hkKG",<br>"token_type": "Bearer",<br>"refresh_token": "8xLOxBtZp8",<br>"expires_in": 3600,<br>"id_token": "a ... b.c ... d.e ... f" // must be a JWT<br>}
        else all users managed by konnectd/ocis
            client->>+konnectd: 2. Client sends the request to the Authorization Server.
            Note over client, konnectd: GET /authorize?<br>response_type=code<br>&scope=openid%20profile%20email<br>&client_id=s6BhdRkqt3<br>&state=af0ifjsldkj<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb HTTP/1.1<br>Host: server.example.com
            Note over user, konnectd: 3. Authorization Server Authenticates the End-User.
            konnectd->>+glauth: LDAP query/bind
            glauth->>+graph: GET user with Basic Auth<br>GraphAPI
            graph->>+accounts: internal GRPC
            accounts-->>-graph: response
            graph-->>-glauth: OData response
            glauth-->>-konnectd: LDAP result
            Note over konnectd,ldap: In case the internal users come from an external ldap they have to be synced to the accounts service to show up as recipients during sharing.

            Note over user, konnectd: 4. Authorization Server obtains End-User Consent/Authorization.
            konnectd-->>-client: 5. Authorization Server sends the End-User back<br>to the Client with an Authorization Code.
            Note over client, konnectd: HTTP/1.1 302 Found<br>Location: https://client.example.org/cb?<br>code=SplxlOBeZQQYbYS6WxSbIA&state=af0ifjsldkj

            client->>+konnectd: 6. Client requests a response using the<br>Authorization Code at the Token Endpoint.
            Note over client, konnectd: POST /token HTTP/1.1<br>Host: server.example.com<br>Content-Type: application/x-www-form-urlencoded<br>grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA<br>&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb
            konnectd-->>-client: 7. Client receives a response that contains an<br>ID Token and Access Token in the response body.
            Note over client, konnectd:  HTTP/1.1 200 OK<br>Content-Type: application/json<br>Cache-Control: no-store<br>Pragma: no-cache<br>{<br>"access_token": "SlAV32hkKG",<br>"token_type": "Bearer",<br>"refresh_token": "8xLOxBtZp8",<br>"expires_in": 3600,<br>"id_token": "a ... b.c ... d.e ... f" // must be a JWT<br>}
        end

        client-->>client: 8. Client validates the ID token and<br>retrieves the End-User's Subject Identifier.

        client->>+proxy: PROPFIND <br> With access token
        proxy-->>-client: 207 Multi-Status
    client-->>-user: List of Files X, Y, Z ...
{{< /mermaid >}}