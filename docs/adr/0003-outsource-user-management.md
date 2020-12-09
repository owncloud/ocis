# 3. Outsource Usermanagement

* Status: accepted <!-- optional -->
* Deciders: @butonic, @mbarz, @kfreitag, @hd, @pmaier1 <!-- optional -->
* Date: 2020-12-09 <!-- optional -->

Technical Story: [Skip account-service by talking to CS3 user-api](https://github.com/owncloud/ocis/pull/1020) <!-- optional -->

## Context and Problem Statement

To attach metadata like shares to users ownCloud relies on persistent, non-reassignable, unique identifiers for users (and files). Email und username can change when a user changes his name. But even the OIDC sub+iss combination may change when the IdP changes. While there is [an account porting protocol](https://openid.net/specs/openid-connect-account-porting-1_0.html) that describes how a relying party such as ownCloud should should behave, it still requires the RP to maintain its own user identifiers.

## Decision Drivers <!-- optional -->

* OCIS should be a single binary that can run out of the box without external dependencies like an LDAP server.
* Time: we want to build a release candiddate asap.

## Considered Options

* Accounts service wraps LDAP
* GLauth wraps accounts service

## Decision Outcome

Chosen option: "Move accounts functionality to GLauth and name it accounts", by moving the existing accounts service file based persistence to glauth and use it as a drop in replacement for an LDAP server. The reverse index and web ui existing in the accounts service will move as well in order to make glauth a standalone, small scale user management with write capabilities.

### Product summary
- GLauth is a drop in user management for small scale deployments.
- OCIS admins can either use the web ui to manage users in glauth or use existing tools in their IDM.
- We hide the complexity by embedding OpenID Provider, an LDAP server and a user management web ui.

### Resulting deployment options
- Single binary: admin can manage users, groups and roles using the built in web ui (glauth)
- External LDAP: OCIS admin needs do use existing tool to manage users
- Separate OCIS and LDAP admin: OCIS admin relies on the LDAP admin to manage users

### Resulting technical implications
- add graphapi to glauth so the ocis web ui can use it to manage users
- make graphapi service to directly talk to an LDAP server so our web ui can use it
- keep the accounts service but embed glauth
- add graph api to the accounts service?

### Positive Consequences <!-- optional -->

* The accounts service (which is our drop in LDAP solution) can be disabled.
* No sync

### Negative Consequences <!-- optional -->

* If users want to store users in their IDM and at the same time guests in a seperate user management we need to implement ldap backends that support more than one LDAP server.

## Pros and Cons of the Options <!-- optional -->

### GLauth wraps accounts service

Currently, the accounts service is the source of truth and we use it to implement user management. <!-- optional -->

* Good, because it solves the problem of storing and looking up an owncloud uuid for a user (and group)
* Good, because we can manage users out of the box
* Good, because we can persist accounts in a CS3 storage provider
* Bad, because it maintains a separate user repository: it needs to either learn or sync users.
* … <!-- numbers of pros and cons can vary -->

### Move accounts functionality to GLauth and name it accounts

We should use an existing ldap server and make GLauth a drop in replacement for it. <!-- optional -->

* Good, because we can use an existing user repository (an LDAP server), no need to sync or learn users.
* Good, because admins can rely on existing user managemen tools.
* Good, because we would have a clear seperation of concerns:
  - users reside in whatever repository, typically an LDAP server
    - could be an existing LDAP server or AD
    - could be our embeddable drop in glauth server
  - we use a service to wrap the LDAP server with other APIs:
    - graph API - ODATA based restful api,
    - [SCIM](http://www.simplecloud.info/) - designed to manage user identities, supported by some IDPs,
    - the current accounts API (which is a protobuf spec following the graph api)
  - our account management ui can use the graph api service which can have different backends
    - an existing ldap server
    - our drop in glauth server (which might serve the graph api itself)
    - the cs3 api + a future guest provisioning api + a future cs3 user provisioning api 
  - all ocis services can use the service registry to look up the accounts service that provides an internal api
    - could be the CS3 user provider (and API)
    - could be the internal protobuf accounts API
  - introduce a new guest provisioning api to CS3 which properly captures our requirement to have them in the user repository
    - guests need to be made available to the firewall
    - storages like eos that integrate with the os for acl based file permissions need a numeric user and group id
* Good, because we can use the CS3 user proviter with the existing ldap / rest driver.
* Bad, because OCIS admins may not have the rights to manage role assignments. (But this is handled at a different department.) 
* Bad, because OCIS admins may not have the rights to disable users if an external LDAP is used instead of the drop in GLauth.
* … <!-- numbers of pros and cons can vary -->
