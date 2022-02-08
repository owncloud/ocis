---
title: "17. Allow read only external User Management"
weight: 17
date: 2022-02-08T10:53:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0017-allow-read-only-external-user-management.md
---

* Status: proposed
* Deciders: @butonic, @micbar, @rhafer
* Date: 2022-02-08

## Context and Problem Statement

oCIS needs to be integrated with various external Authentication and Identity Management Systems.

Sidenote: There is a difference between users, identities and accounts: A user may have multiple
identities whith which he can authenticate, e.g. his facebook, twitter, microsoft or google
identity. Multiple identities can be linked to an account in ocis, allowing to fall back to another
identity provider should one of them shut down. This also allows migrating from one identity
provider to another.

There are different cases where oCIS requires access to users:

1. While we settled on using OpenID Connect (OIDC) as the authentication protocol for oCIS, we
   we need to build a user object during authentication with at least an account UUID (to identify
   the account) and the email or a name (for display purposes). 
2. When searching for share recipients we need to be able to query existing users in the external
   identity management system
3. When listing files we need to be able to look up a users display properties (username, email,
   avatar...) based on the account UUID

oCIS internally relies on a stable and persistent identifier (e.g. a UUID) for accounts in order to
implement permissions and sharing. Unfortunately, some deployments are unable to deliver this kind
of stable identifier for users:

- In OIDC itself the only stable identifier that is guaranteed to be provided by the IDP is
  combination of the sub and iss claims. IDPs can optionally return other claims, but we cannot
  rely on a specific claim being present.
- When no other services (LDAP, SCIM, ...) is available that could be used look up a user UUID


## Decision Drivers

* oCIS should be a single binary that can run out of the box without external dependencies like an
  LDAP server.
* Time: we want to build a release candidate asap.
* oCIS should be easy to integrate with standard external identity mangement systems

## Considered Options

* External identity management system is writeable and has all necessary APIs
* External identity management system is read only and provides an interface to query users
* External identity management system is read only and does NOT provide an API to query users

## Decision Outcome

tbd

### Positive Consequences: <!-- optional -->


### Negative consequences: <!-- optional -->


## Pros and Cons of the Options <!-- optional -->

### External identity management system is writeable and has all necessary APIs

IdP sends all necessary claims: uuid, username, email, displayname, avatar url IdP allows lookup of
display properties by the uuid or email/username

* Good, because we can fully rely on the external identity management system
* Bad, because we need write access to provision guest accounts (very few customers are willing to
  provide that)

### External identity management system is read only and provides an interface to query users (e.g. Coporate Active Directy)

IdP sends sub & iss and mail or username claims, Identity Management System provides APIs (e.g.
LDAP, SCIM, REST ...) to lookup additional user information. All services use the CS3 API to look up
the account for the given email or username, where CS3 then uses a backend that relies on the APIs
provided by the IdM.

* Good, because we can rely on the external identity management
* Good, because ocis services only need to know about the CS3 user provider API, which acts as an
  abstraction layer for different identitiy management systems
* Good, because there is only an single source of truth (the external IdM) and we don't need to
  implement a synchronization mechanism to maintain an internal user database (we will likely need
  some form of caching though, see below)
* Bad, because the identity managment needs to provide a stable, persistent, non-reassignable user
  identifier for an account, e.g. `owncloudUUID` or `ms-DS-ConsistencyGuid`
* Bad, because we need to implment tools that can change the account id when it did change anyway
* Bad, because without caching we will hammer the identity management system with lookup requests

### External identity management system is read only and does NOT provide an API to query users

Idp sends sub & iss and mail or username claims. We need to provision an internal account mapping
upon first login of a user to be able to look up user properties by account id.

* Good, because this has very little external requirements
* Good, because we have accounts fully under our control
* Bad, because we have to provide the user lookup APIs
* Bad, because users will only a visible after the first login

## Links <!-- optional -->

* [Link type] [Link to ADR] <!-- example: Refined by [ADR-0005](0005-example.md) -->
* … <!-- numbers of links can vary -->
* supersedes [3. Use external User Management]({{< ref "0003-external-user-management.md" >}})
