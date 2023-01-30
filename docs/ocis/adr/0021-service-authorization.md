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

There are three levels of security checks in a microservice web application that uses openid connect:
1. **scope claims** limit the possible operations to what the user (or admin on behalf of the organization) consented to
2. **service authorization** limit the possible operations to what specific services are allowed to do, on behalf of users or even without them
3. **permission checks** limit the possible operations to the relationships between subject, permission and resource allow

This ADR deals with service authorization. 

Some services need access to file content without a user being logged in. We currently pass the owner or manager
of a space in events which allows tśhe search service to impersonate that user to extract metadata from the changed resource.
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

### Consequences

* Good, because {positive consequence, e.g., improvement of one or more desired qualities, …}
* Neutral, because {neutral consequence, e.g., compromising one or more desired qualities, …}
* Bad, because {negative consequence, e.g., compromising one or more desired qualities, …}

## Pros and Cons of the Options

### Service Accounts

Implement service accounts to mimic the MS graph servicePrincipals for an application. The appRoles can be used to limit the interactions to other services, e.g. `Files.Read.All` to allow the servicePrincipal reading all files in all drives of a tenant or [`User.ReadWrite.All`](https://learn.microsoft.com/en-us/graph/permissions-reference#user-permissions) to manage users in an organization. 
An application is represented by a dedicated servicePrincipal per tenant, which is created when granting an application access to a tenant (currently instace) by an admin. These servicePrincipals then get an appRoleAssignment for the appRoles defined by the application.
An application could be the search service (or only a content extraction service), a workflow service, a publication service or any other service.

#### Q&A

* Which services need a service account?
Every service. Either on behalf of a logged in user (thumbnails or graph), or without a user being logged in (antivirus, search indexing). The MS graph api calls the former *delegated permissions* and the latter *application permissions*. See the [MS Graph permissions overview](https://learn.microsoft.com/en-us/graph/permissions-overview?tabs=http) 

Service | Delegated Permissions | Application Permissions | Comment
-|-|-
graph | `Files.ReadWrite.All`, `Drive.Read` | | The graph translates all requests of the currently logged in user into CS3 requests. For admins and space admins `Drives.ReadWrite.All` would allow managing drives.
search | `Files.Read.All` | | Search needs to be able to read the indexed content, which is equivalent to all files, the currently logged in user has access to
indexer | | `Files.Read.All` | The indexing process is triggered asynchronously and the admin needs to consent on behalf of all users / the organization / tenant to allow indexing all drives. Since indexing is currently part of the search service we should split the two services to prevent users from gaining access to an indexer process with `Files.Read.All` application permissions by exploiting a bug in the search process.
thumbnail | `Files.Read.All` | | The thumbnailer needs to access file content the user has access to, to generate a thumbnail. It should store them on a service storage (which is what the go micro store interface should be used for). If we want to pre generate thumbnails the service would need the `Files.Read.All` application permission, but again, a dedicated service would be preferable for security reasons.
antivirus | | `Files.Read.All` | As part of postprocessing this is event triggered and requires an application permission.


#### Why don't we need service accounts for services that manage and list users?
Actually, we do. One side of introducing service accounts is giving every service a service account. The other is actually making called services check the permissions of the service account. We first need service accounts for search indexing and antivirus because they actually solve a problem there. Then we can gradually increase the security of the inter service communication by requiring more and more permissions, e.g. for managing user accounts.

#### Where do we store them? In the same user backend?
In the CS3 API we can find USER_TYPE_SERVICE and USER_TYPE_APPLICATION, where application is something like collabora and service is afaict sth like service accounts.

#### Do we have to implement a new reva auth-service manager?

#### Do graph permission names have to bleed into the internal CS3 implementation?
Maybe not, AFAICT we could internally use spicedb and use it to implement the permission checks.

#### Why can't we just mint a token and bake the permissions the service has into it.
For normal users this is called the scope and it covers the permissions the use has consented to. When an access token arrives at the proxy we can use the scope claim which is filled by the IdP to limit the request to that scope. However, we will not be able to teach every IdP every scope that users should be able to consent to. It is far more practical to read a role claim and mint the scopes based on that, e.g. a role user would get the `Files.ReadWrite.All` permission that scopes the request to read all his files. A service that uses delegated permission can during this request only access the users files and all files the user has access to. We can add a scope middleware to the graph api (or other http endpoints) that checks if the user has consonted to the scope required for the request.

This does not free us from the actual permission check to determine if the user can actually read or write the file. While `Files.ReadWrite.All` means that the service could write to a file that was shared with the logged in user an actual permission check needs to determine if the file was shared with write permission or if it was read only.

#### How do we check the service account has the permission to make a call?
In the same way we currently check permissions. We delegate that to the permissions service. Hm but that means we need to make two permission checks: one for the service account and one for the logged in user, if present. We need to determine if the service is allowed to make a request and if the user is allowed to make a request. We want to protect against .g. a compromised thumbnail service that should only make read requests but then tries to write a file. This needs to be checked at the storage provider.

What if we logically AND the delegated permission with the users scope?

How does go micro solve this?

To be clear: this has nothing to do with the scope.


#### How do we "provision" service accounts?

#### How/where do we assign roles/permissions to them?

 
How do services declare what permissions they need?
  Especially useful when the admin wants to use a (maybe closed source) 3rd party extension. (Some day in the future)

* Good, because we could replace machine auth with specific service accounts and no longer have to distribude a shared secret everywhere
* Bad, because we don't know if a there are places in the code that try to look up a user with USER_TYPE_SERVICE at the cs3 users service ... they might not exist there ... or do we have to implement a userregistry, similar to the authregistry?
* Bad, because we have to provision and manage service accounts on init
* Bad, because we have to write codemanage service accounts in the admin ui


### Auth Manager type Space-Owner

We could implement a new auth manager that can authenticate space owners, a CS3 user type we introduced for project spaces which 'have no owner', only one or more managers.

* Good, because it reuses the space owner user type
* Bad, because the space owner always has write permisson
* Bad, because we don't know if a there are places in the code that try to look up a user with USER_TYPE_SPACE_OWNER at the cs3 users service ... they might not exist there ... or do we have to implement a userregistry, similar to the authregistry?
* Bad, because it feels like another hack and does not protect against compromized services that try to execute operations that the user did not consent to.

## Links

* [MS Graph servicePrincipal](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-1.0)
* [reva auth managers](https://reva.link/docs/config/packages/auth/manager/) - lacks docs for `auth_machine`, to be found [in the code](https://github.com/cs3org/reva/blob/edge/pkg/auth/manager/machine/machine.go)
