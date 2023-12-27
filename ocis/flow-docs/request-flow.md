---
title: "Request Flow"
date: 2020-04-27T16:07:00+01:00
weight: 45
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/flow-docs
geekdocFilePath: request-flow.md
---


## Request Flow

The following sequence diagram describes the general request flow. It shows where account provisioning and token minting are happening:

{{< mermaid class="text-center">}}
sequenceDiagram
    %% we have comments!! \o/
    participant user as User
    participant client as Client
    participant proxy as ocis-proxy
    participant idp as IdP
    participant accounts as ocis-accounts
    participant ldap as corporate LDAP server

    user->>+client: What is the content of my home?

        client->>+proxy: PROPFIND <br> Bearer auth using oidc auth token
        Note over client,proxy: What is in a bearer token? <br> The spec recommends opaque tokens. <br> Treat it as random byte noise.
        Note over client,proxy: the proxy MUST authenticate users <br> using ocis-accounts because it needs <br> to decide where to send the request
        %% Mention introspection endpoint for opaque tokens
        %% idp uses jwt, so we can save a request
        %% either way the token can be used to look up the sub and iss of the user

            %% or is token check enough?
            proxy->>+idp: GET /userinfo
            alt userinfo succeeds

                idp-->>proxy:  200 OK
                Note over proxy,accounts: Content-Type: application/json<br>{<br>"sub": "248289761001",<br>"name": "Jane Doe",<br>"given_name": "Jane",<br>"family_name": "Doe",<br>"preferred_username": "j.doe",<br>"email": "janedoe@example.com",<br>"picture": "http://example.com/janedoe/me.jpg"<br>}
                %% see: https://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse

            else userinfo fails

                idp-->>-proxy: 401 Unauthorized
                Note over proxy,accounts: WWW-Authenticate: error="invalid_token",<br>error_description="The Access Token expired"

        proxy-->>client: 401 Unauthorized or <br>302 Found with redirect to idp
        Note over client: start at login flow<br> or refresh the token

            end

            proxy->>+accounts: TODO API call to exchange sub@iss with account UUID
            Note over proxy,accounts: does not autoprovision users. They are explicitly provisioned later.

            alt account exists or has been migrated

                accounts-->>proxy: existing account UUID
            else account does not exist

                opt oc10 endpoint is configured
                Note over proxy,oc10: Check if user exists in oc10
                    proxy->>+oc10: GET /apps/graphapi/v1.0/users/&lt;uuid&gt;
                    opt user exists in oc10
                        oc10-->>-proxy: 200
                        %% TODO auth using internal token
                        proxy->>+oc10: PROPFIND
                        Note over proxy,oc10: forward existing bearer auth
                        oc10-->>-proxy: Multistatus response
            proxy-->>client: Multistatus response
    client-->>user: List of Files X, Y, Z ...
                    end
                end

                Note over proxy,accounts: provision a new account including displayname, email and sub@iss <br> TODO only if the user is allowed to login, based on group <br> membership in the ldap server
                proxy->>proxy: generate new uuid
                proxy->>+accounts: TODO create account with new generated uuid
                accounts-->>-proxy: OK / error

            else account has been disabled

                accounts-->>-proxy: account is disabled
        proxy-->>client: 401 Unauthorized or <br>302 Found with redirect to idp
        Note over client: start at login flow<br> or refresh the token

            end
            proxy->>proxy: store uuid in context

            %% what if oc10 does not support a certain request / API

            proxy->>proxy: mint an internal jwt that includes the UUID and username using revas `x-access-token` header
            proxy->>+reva: PROPFIND <br>Token auth using internal JWT
            reva-->>-proxy: Multistatus response
        proxy-->>-client: Multistatus response

    client-->>-user: List of Files X, Y, Z ...
{{< /mermaid >}}
