---
title: "17. Allow read only external User Management"
---

* Status: proposed
* Deciders: @butonic, @micbar, @rhafer
* Date: 2022-02-08

## Context and Problem Statement

oCIS needs to be integrated with various external Authentication and Identity Management System. We
settled on Open ID Connect (OIDC) as the central authentication protocol for OCIS.

OCIS internally relies on a stable and persistent identifier (e.g. a UUID) for accounts in order to
implement permissions and sharing. Unfortunately, some deployments are unable to deliver this kind
of stable identifier for users:

- In OIDC itself the only stable identifier that is guaranteed to be provided by the IDP is
  combination of the sub and iss claims. IDPs can optionally return other claims, but we might not
  be able to rely on a specific claim being present.
- When no other services (LDAP, SCIM, ...) is available that could be used look up a user UUID

Furthermore, there is a difference between users, identities and accounts: A user may have multiple
identities whith which he can authenticate, e.g. his facebook, twitter, microsoft or google
identity. Multiple identities can be linked to an account in ocis, allowing to fall back to another
identity provider should one of them shut down. This also allows migrating from one identity
provider to another.

There are three cases that require access to users:a

1. During authentication we neet to build a user object with at least an account uuid (to identify
   the account) and the email (for display purposes)
2. When searching for recipients we need to be able to query existing users in the external identity
   management system
3. When listing files we need to be able to look up a users display properties (username, email,
   avatar...) based on the account uuid

## Decision Drivers

* oCIS should be a single binary that can run out of the box without external dependencies like an
  LDAP server.
* Time: we want to build a release candidate asap.
* oCIS should be able to be easily integrated with external standard identity mangement systems

## Considered Options

* External identity management system is writeable and has all necessary APIs
* External identity management system is read only and provides an interface to query users (e.g. 
* IdP is read only and does not provide an API to query users

## Decision Outcome

tbd

### Positive Consequences: <!-- optional -->


### Negative consequences: <!-- optional -->


## Pros and Cons of the Options <!-- optional -->

### External identity management system is writeable and has all necessary APIs

IdP sends all necessary claims: uuid, username, email, displayname, avatar url IdP allows lookup of
display properties by the uuid or email/username

* Good, because we can fully rely on the external identity management system
* Bad, because we need write access to provision guest accounts

### External identity management system is read only and provides an interface to query users (e.g. Coporate Active Directy)

IdP ends sub & iss and mail or username claims, Identity Management System provides Interfaces (e.g.
LDAP) to lookup additional user information. All services use the CS3 API to look up the account for
the given email or username, where CS3 provides backends for LDAP, SCIM, REST ...

* Good, because we can rely on the external identity management
* Bad, because the Identity managment needs to provide a stable, persistent, non-reussignable user
  identifier to identify the account, e.g. owncloudUUID or ms-DS-ConsistencyGuid
* Bad, because we need to implment tools that can change the account id when it did change anyway
* Bad, because we will hammer the identity management system with lookup requests (can mostly be
  mitigated with caching)

### IdP is read only and does not provide an API to query users

Idp sends sub & iss and mail or username claims. We need to provision an internal account mapping to
look up user properties by account id.

* Good, because this has very little external requirements
* Bad, because we have to provide the user lookup APIs

## Links <!-- optional -->

* [Link type] [Link to ADR] <!-- example: Refined by [ADR-0005](0005-example.md) -->
* … <!-- numbers of links can vary -->
* supersedes [3. Use external User Management]({{< ref "0003-external-user-management.md" >}})
