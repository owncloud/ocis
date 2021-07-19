---
title: "1. Introduce an accounts service"
weight: 1
date: 2020-06-15T20:21:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0001-introduce-accounts-service.md
---

* Status: superseded by [ADR-0003]({{< ref "0003-external-user-management.md" >}})
* Deciders: @butonic, @felixboehm, @micbar, @pmaier1
* Date: [2020-06-15](https://github.com/owncloud/ocis-accounts/pull/34/commits/2fd05e2b6fe2a47c687bd0c0bc5e1b5c48a585b2)

Technical Story: [persist accounts](https://github.com/owncloud/ocis-accounts/pull/34)

## Context and Problem Statement

To attach metadata like shares to users ownCloud relies on persistent, non-reassignable, unique identifiers for users (and files). Email und username can change when a user changes his name. But even the OIDC sub+iss combination may change when the IdP changes. While there is [an account porting protocol](https://openid.net/specs/openid-connect-account-porting-1_0.html) that describes how a relying party (RP) such as ownCloud should behave, it still requires the RP to maintain its own user identifiers.

## Decision Drivers

* oCIS should be a single binary that can run out of the box without external dependencies like an LDAP server.
* Time: we want to build a release candidate asap.
* Firewalls need access to guests, typically via LDAP.
* Not all external LDAPs are writeable for us to provision Guest accounts.
* We see multiple LDAP servers in deployments. Being able to handle them is important and should be covered by using OIDC + being able to query multiple LDAP servers.

## Considered Options

* Accounts service wraps LDAP
* [GLAuth](https://github.com/glauth/glauth) wraps accounts service

## Decision Outcome

Chosen option: "GLAuth wraps accounts service", because we need write access to provision guest accounts and GLAuth currently has no write support.

### Positive Consequences

* We can build a self contained user management in the accounts service and can adjust it to our requirements.
* We do not rely on an LDAP server which would only be possible by implementing write support in the LDAP libraries used by GLAuth (hard to estimate effort, when will that be merged upstream).

### Negative Consequences

* We need to spend time on implementing user management

## Pros and Cons of the Options

### Accounts service wraps LDAP

* Bad, because not all external LDAPs are writeable for us to provision Guest accounts.
