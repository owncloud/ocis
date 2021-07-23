---
title: "3. Use external User Management"
weight: 3
date: 2020-12-09T20:21:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0003-external-user-management.md
---

* Status: accepted
* Deciders: @butonic, @micbar, @dragotin, @hodyroff, @pmaier1
* Date: 2020-12-09

Technical Story: [Skip account-service by talking to CS3 user-api](https://github.com/owncloud/ocis/pull/1020)

## Context and Problem Statement

To attach metadata like shares to users ownCloud relies on persistent, non-reassignable, unique identifiers for users (and files). Email and username can change when a user changes his name. But even the OIDC sub+iss combination may change when the IdP changes. While there is [an account porting protocol](https://openid.net/specs/openid-connect-account-porting-1_0.html) that describes how a relying party (RP) such as ownCloud should behave, it still requires the RP to maintain its own user identifiers.

## Decision Drivers

* oCIS should be a single binary that can run out of the box without external dependencies like an LDAP server.
* Time: we want to build a release candidate asap.
* oCIS should be able to be easily integrated with standard user management components

## Considered Options

* Accounts service wraps LDAP
* [GLAuth](https://github.com/glauth/glauth) wraps accounts service

## Decision Outcome

Chosen option: "Move accounts functionality to GLAuth and name it accounts", by moving the existing accounts service file based persistence to GLAuth and use it as a drop in replacement for an LDAP server. The reverse index and web UI existing in the accounts service will move as well in order to make GLAuth a standalone, small scale user management with write capabilities.

### Product summary
- GLAuth is a drop in user management for small scale deployments that do not rely on an actual LDAP server.
- oCIS admins can either use the web UI to manage users in GLAuth or use existing tools in their IDM.
- We hide the complexity by embedding an OpenID Provider, an LDAP server and a user management web UI.

### Resulting deployment options
- Use internal user management
  - Recommended for small scale use cases and simple deployments
  - Users, groups and roles are stored and managed within GLAuth
- Use external user management
  - Recommended for mid and large scale use cases
  - Users, groups and roles are stored and managed within an external LDAP / AD / IDM
  - Separate oCIS and LDAP admin: oCIS admin relies on the LDAP admin to manage users
- User permissions for roles are always managed in oCIS (settings service) because they are specific to oCIS

### Resulting technical implications
- Make the file based reverse index a standalone library
- Contribute to GLAuth
  - Add ms graph based rest API to manage users, groups and roles (the LDAP lib is currently readonly)
  - Add web UI to glauth that uses the ms graph based rest API to manage users
  - Add a backend that uses the file based reverse index, currently living in the oCIS accounts service
  - Move fallback mechanism from ocis/glauth service to upstream GLAuth to support multiple LDAP servers
    - Make it a chain to support more than two LDAP servers
    - Document the implications for merging result sets when searching for recipients
    - At least one writeable backend is needed to support creating guest accounts
- Make all services currently using the accounts service talk to the CS3 userprovider
- To support multiple LDAP servers we need to move the fallback mechanism in ocis/glauth service to upstream GLAuth
- The current CS3 API for user management should be enriched with pagination, field mask and a query language as needed
- properly register an [auxiliary LDAP schema that adds an ownCloudUUID attribute to users and groups](https://github.com/owncloud/ocis/blob/c8668e8cb171860c70fec29e5ae945bca44f1fb7/deployments/examples/cs3_users_ocis/config/ldap/ldif/10_owncloud_schema.ldif)

### Positive Consequences

* The accounts service (which is our drop in LDAP solution) can be dropped. The CS3 userprovider service becomes the only service dealing with users.
* No sync

### Negative Consequences

* If users want to store users in their IDM and at the same time guests in a seperate user management we need to implement GLAuth backends that support more than one LDAP server.

## Pros and Cons of the Options

### GLAuth wraps accounts service

Currently, the accounts service is the source of truth and we use it to implement user management.

* Good, because it solves the problem of storing and looking up an owncloud UUID for a user (and group)
* Good, because we can manage users out of the box
* Good, because we can persist accounts in a CS3 storage provider
* Bad, because it maintains a separate user repository: it needs to either learn or sync users.

### Move accounts functionality to GLAuth and name it accounts

We should use an existing LDAP server and make GLAuth a drop in replacement for it.

* Good, because we can use an existing user repository (an LDAP server), no need to sync or learn users.
* Good, because admins can rely on existing user management tools.
* Good, because we would have a clear separation of concerns:
  - users reside in whatever repository, typically an LDAP server
    - could be an existing LDAP server or AD
    - could be our embeddable drop in glauth server
  - we use a service to wrap the LDAP server with other APIs:
    - ms graph API - ODATA based restful API,
    - [SCIM](http://www.simplecloud.info/) - designed to manage user identities, supported by some IDPs,
    - the current accounts API (which is a protobuf spec following the ms graph API)
  - our account management UI can use the ms graph based API service which can have different backends
    - an existing LDAP server
    - our drop in glauth server (which might serve the ms graph based API itself)
    - the CS3 API + a future guest provisioning API + a future CS3 user provisioning API (or [generic space provisioning](https://github.com/cs3org/cs3apis/pull/95))
  - all oCIS services can use the service registry to look up the accounts service that provides an internal API
    - could be the CS3 user provider (and API)
    - could be the internal protobuf accounts API
  - introduce a new guest provisioning API to CS3 which properly captures our requirement to have them in the user repository
    - guests need to be made available to the firewall
    - storages like EOS that integrate with the os for acl based file permissions need a numeric user and group id
* Good, because we can use the CS3 user provider with the existing LDAP / rest driver.
* Bad, because oCIS admins may not have the rights to manage role assignments. (But this is handled at a different department.) 
* Bad, because oCIS admins may not have the rights to disable users if an external LDAP is used instead of the drop in GLAuth.

## Links
* supersedes [ADR-0001]({{< ref "0001-introduce-accounts-service.md" >}})
