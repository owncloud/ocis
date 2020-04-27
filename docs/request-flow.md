---
title: "Request Flow"
date: 2020-04-27T16:07:00+01:00
weight: 45
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: request-flow.md
---


## Request Flow

The following sequence diagram describes the general request flow:

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
        Note right of client: What is in a bearer token? <br> The spec recommends opaque tokens. <br> So it is just random byte noise.
        %% Mention introspection endpoint for opaque tokens
        %% konnectd uses jwt, so we can save a request
        %% either way the token can be used to look up the sub and iss of the user

            %% or is token check enough?
            proxy->>+idp: GET /userinfo
            idp-->>-proxy: JSON response
            Note right of proxy: the result contains <br> the sub of the user
            %% see: https://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse

            proxy->>+accounts: TODO API call to exchange sub@iss with account UUID

                alt internal account
                    accounts->>+ldap: is user allowed to use ocis
                    ldap-->>-accounts: yes/no - group based
                else guest account
                    accounts->>accounts: check if is valid guest account
                end


            accounts-->>-proxy: new or existing account UUID / error
            Note right of accounts: actually this provisions <br> the account including <br> displayname, email and <br> sub@iss if the user is <br> allowed to login, based <br> on group membership <br> in the ldap server


            Note right of proxy: the proxy MUST <br> authenticate users <br> using ocis-accounts <br> because it needs to <br> decide where to <br> send the request

            Note right of proxy: forward request to <br> ocis or oc10
            %% what if oc10 does not support a certain request / API
            alt user is migrated

                Note right of proxy: mint an internal jwt <br> token that includes <br> the UUID and username
                proxy->>+reva: PROPFIND <br> Bearer auth using internal JWT
                reva-->>-proxy: Multistatus response

            else user is not migrated

                Note right of proxy: forward existing bearer auth?
                proxy->>+oc10: PROPFIND <br> Bearer auth using internal JWT
                %% TODO auth using internal token?
                oc10-->>-proxy: Multistatus response

            end


        proxy-->>-client: Multistatus response

    client-->>-user: List of Files X, Y, Z ...
{{< /mermaid >}}