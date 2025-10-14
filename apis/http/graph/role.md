---
title: Role
weight: 60
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/http/graph
geekdocFilePath: permissions.md
---

{{< toc >}}

## Role API

The Roles API is implementing a subset of the functionality of the
[MS Graph Role Management](https://learn.microsoft.com/en-us/graph/api/resources/rolemanagement?view=graph-rest-1.0).

## Role Management

### List roleDefinitions `GET /v1beta1/roleManagement/permissions/roleDefinitions`

https://owncloud.dev/libre-graph-api/#/roleManagement/ListPermissionRoleDefinitions

### Get unifiedRoleDefinition `GET /drives/{drive-id}/items/{item-id}/permissions/{perm-id}`

https://owncloud.dev/libre-graph-api/#/roleManagement/GetPermissionRoleDefinition

## Role Assignment

### Get appRoleAssignments of a user `GET /v1.0/users/{user-id}/appRoleAssignments`

https://owncloud.dev/libre-graph-api/#/user.appRoleAssignment/user.ListAppRoleAssignments

### Grant an appRoleAssignment to a user `POST /v1.0/users/{user-id}/appRoleAssignments`

https://owncloud.dev/libre-graph-api/#/user.appRoleAssignment/user.CreateAppRoleAssignments

### Delete the appRoleAssignment from a user `DELETE /v1.0/users/{user-id}/appRoleAssignments/{appRoleAssignment-id}`

https://owncloud.dev/libre-graph-api/#/user.appRoleAssignment/user.DeleteAppRoleAssignments