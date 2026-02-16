---
title: "21. Service accounts"
date: 2023-01-18T16:07:00+01:00
weight: 21
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0021-service-accounts.md
---

* Status: proposed
* Deciders: [@butonic](https://github.com/butonic), [@c0rby](https://github.com/c0rby)
* Date: 2023-01-18

## Context and Problem Statement

There are three levels of security checks in a microservice web application that uses OpenID Connect:
1. **scope claims** limit the possible operations to what the user (or admin on behalf of the organization) consented to
2. **service authorization** limit the possible operations to what specific services are allowed to do, on behalf of users or even without them
3. **permission checks** limit the possible operations to the relationships between subject, permission and resource allow

This ADR deals with a prerequisite for service authorization: service accounts.

Some services need access to file content without a user being logged in. We currently pass the owner or manager
of a space in events which allows the search service to impersonate that user to extract metadata from the changed resource.
There are two problems with this:
1. The service could get all permissions of the user and gain write permission
2. There is a race condition where the user in the event might no longer have read permission, causing the index to go stale

The race condition will become more of an issue when we start working on a workflow engine.

How can we grant services the least amount of permissions required for their purpose?

## Decision Drivers

* It should be possible to represent this as servicePrincipals in the libregraph API, similar to the [MS Graph servicePrincipal](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-1.0).
* Services should check permissions using the oCIS permissions or reva auth service, we don't want to introduce a new mechanism for this

## Considered Options

* [Service Accounts](#service-accounts)
* [Impersonate Space-Owners](#impersonate-space-owners)

## Decision Outcome

Chosen option: [Service Accounts](#service-accounts)

### Consequences

* Good, because it allows provisioning permissions for services
* Good, because it uses existing CS3 concepts
* Good, because it uses the existing permissions service
* Good, because it can be mapped to libre graph permissions
* Bad, because we have to make the reva auth manager aware of CS3 [`USER_TYPE_SERVICE`](https://cs3org.github.io/cs3apis/#cs3.identity.user.v1beta1.UserType)
* Bad, because we have to provision and manage service accounts on init
* Bad, because external APIs may need to filter out service accounts
* Bad, because we need to persist service accounts in addition to normal user accounts

## Pros and Cons of the Options

### Service Accounts

Make the reva auth manager and registry aware of CS3 users of type [`USER_TYPE_SERVICE`](https://cs3org.github.io/cs3apis/#cs3.identity.user.v1beta1.UserType). Then we can provision service accounts at oCIS initialization and use the permissions service to check permissions.
When assigning permissions we use the permission constraints to define the scope of permissions, see [Permission Checks](#permission-checks) for more details.

To authenticate service accounts the static reva auth registry needs to be configured with a new auth provider for type `service`. The actual provider can use a plain JSON file or JSONCS3 that is provisioned once with `ocis init`. TODO Furthermore, the user provider needs to be able to return users for service accounts.


* Good, because we could replace machine auth with specific service accounts and no longer have to distribute a shared secret everywhere
* Bad, because we don't know if a there are places in the code that try to look up a user with USER_TYPE_SERVICE at the cs3 users service ... they might not exist there ... or do we have to implement a userregistry, similar to the authregistry?
* Bad, because we have to provision and manage service accounts on init
* Bad, because we have to write code to manage service accounts or at least filter them out in the admin ui


### Impersonate Space-Owners

We could implement a new auth manager that can authenticate space owners, a CS3 user type we introduced for project spaces which 'have no owner', only one or more managers.

* Good, because it reuses the space owner user type
* Bad, because the space owner always has write permission
* Bad, because we don't know if a there are places in the code that try to look up a user with USER_TYPE_SPACE_OWNER at the cs3 users service ... they might not exist there ... or do we have to implement a userregistry, similar to the authregistry?
* Bad, because it feels like another hack and does not protect against compromised services that try to execute operations that the user did not consent to.

## Links

* [MS Graph servicePrincipal](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-1.0)
* [reva auth managers](https://reva.link/docs/config/packages/auth/manager/) - lacks docs for `auth_machine`, to be found [in the code](https://github.com/cs3org/reva/blob/edge/pkg/auth/manager/machine/machine.go)

## Permission checks
When checking permissions we do not check for global permissions but for the concrete permission. Global permissions describe permissions that are used when assigning permissions, e.g. the index service account has the read permission constrained to tenant. The concrete permission check always contains a resource and a specific permission like `Resource.Read` or `Space.Delete`. That we currently check if a user has the `delete-all-spaces` permission is wrong. It should instead check if the user has the permission `Space.Delete` on a specific space. The permissions service can implement the check by taking the permission constraint into account.

Another example would be a `Resource.Read` check for a specific resource. Normal users like the demo users Einstein and Marie would have the permission `Resource.ReadWrite` with the constraint ALL (which limits them to all files they own and that have been shared with them). The permissions service can return true. Service accounts like the indexer would have  `Resource.Read` with the constraint TENANT and thus be granted read access to all resources.

In the storage drive implementation we can check the ACLs first (which would allow service accounts that are known to the underlying storage system, e.g. EOS to access the resource) and then make a call to the permissions service. At least for the Read Resource permission. Other permission checks can be introduced as needed.

The permission names and constraints are different from the MS Graph API. Giving permission like [`Files.ReadWrite.All`](https://learn.microsoft.com/en-us/graph/permissions-reference#user-permissions) a different meaning, depending on the type of user (for normal users it means all files they have access to, for service accounts it means all files in the organization) is a source of confusion which only gets worse when there are two different UUIDs for this.
