---
title: "21. Service authorization"
date: 2023-01-18T16:07:00+01:00
weight: 21
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0020-service-authorization.md
---

* Status: proposed
* Deciders: @butonic
* Date: 2023-01-18

## Context and Problem Statement

Some services need access to file content without a user being logged in. We currently pass the owner or manager
of a space in events which allows the search service to impersonate that user to extract metadata from the changed resource.
There are two problems with this:
1. The service could get all permissions of the user and gain write permission
2. There is a race condition where the user in the event might no longer have read permission, causing the index to go stale

The race condition will become more of an issue when we start working on a workflow engine.

How can we grant services the least amount of permissions required for their purpose?

## Decision Drivers

* It should be possible to represent this as servicePrincipals in the libregraph API, similar to the [MS Graph servicePrincipal](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-1.0).
* Services should check permissions using the ocis permissions or reva auth service, we don't want to introduce a new mechanism for this

## Considered Options

* [Service Accounts](#service-accounts)
* [Auth Manager type Space-Owner](#auth-manager-type-space-owner)

## Decision Outcome

Chosen option: ?

### Positive Consequences

* ???
* ???
* ???

### Negative consequences

* ???
* ???
* ???

## Pros and Cons of the Options

### Service Accounts

Implement service accounts to mimic the MS graph servicePrincipals for an application. The appRoles can be used to limit the interactions to other services, e.g. `Files.Read.All` to allow the servicePrincipal reading all files in all drives of a tenant. 
An application is represented by a dedicated servicePrincipal per tenant, which is created when granting an application access to a tenant (currently instace) by an admin. These servicePrincipals then get an appRoleAssignment for the appRoles defined by the application.
An application could be the search service (or only a content extraction service), a workflow service, a publication service or any other service.

TODO which services need a service account?
TODO do we store them in the same user backend? CS3 has USER_TYPE_SERVICE and USER_TYPE_APPLICATION, where application is sth like collabora and service is afaict sth like service accounts. No, but we have to implement a new reva auth-service manager.

* Good, because we could replace machine auth with specific service accounts and no longer have to distribude a shared secret everywhere
* Bad, because we don't know if a there are places in the code that try to look up a user with USER_TYPE_SERVICE at the cs3 users service ... they might not exist there ... or do we have to implement a userregistry, similar to the authregistry?
* Bad, because we have to provision and manage service accounts on init
* Bad, because we have to write codemanage service accounts in the admin ui


### Auth Manager type Space-Owner

We could implement a new auth manager that can authenticate space owners, a CS3 user type we introduced for project spaces which 'have no owner', only one or more managers.

* Good, because it reuses the space owner user type
* Bad, because the space owner always has write permisson
* Bad, because we don't know if a there are places in the code that try to look up a user with USER_TYPE_SPACE_OWNER at the cs3 users service ... they might not exist there ... or do we have to implement a userregistry, similar to the authregistry?

## Links

* [MS Graph servicePrincipal](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-1.0)
* [reva auth managers](https://reva.link/docs/config/packages/auth/manager/) - lacks docs for `auth_machine`, to be found [in the code](https://github.com/cs3org/reva/blob/edge/pkg/auth/manager/machine/machine.go)
