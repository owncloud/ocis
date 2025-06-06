---
title: Unified Roles
date: 2025-06-06T00:00:00+00:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/graph
geekdocFilePath: unified-roles.md
geekdocCollapseSection: true
---

{{< toc >}}

## Create a New Built-in Role

To create a new built-in role, it is necessary to:

- Create role in Reva
- Update Reva in oCIS
- Add newly created role to oCIS

### Add Role to Reva

In the [Reva repository](https://github.com/owncloud/reva), add the role into the `/pkg/conversions/role.go` file:

1. Add a role name constant for this role. See the existing ones for how this is setup
1. Add a new function to create a new role struct with the role name constant and desired permissions
1. Add the role to the `RoleFromName` function
1. In `/pkg/conversions/role_test.go`, extend unit tests to cover the new role

### Add Role to oCIS

After adding the role to Reva and updating the Reva in oCIS, it is necessary to add the role to oCIS as well:

1. Generate UUID for the role
1. In `/services/graph/pkg/unifiedrole/roles.go`, add the role ID generated in the first step as a constant
1. In `/services/graph/pkg/unifiedrole/roles.go`, add translatable role display name and description variables
1. In `/services/graph/pkg/unifiedrole/roles.go`, add role variable with a function to create the role\
The function should first create the role struct using the function from Reva and return the `UnifiedRoleDefinition` struct.
1. In `/services/graph/pkg/unifiedrole/filter.go`, add the role into the `buildInRoles` function
1. In `/services/graph/pkg/config/defaults/defaultconfig.go`, if the role is not intended to be enabled by default, add the role ID into the `_disabledByDefaultUnifiedRoleRoleIDs` constant
1. In `/services/graph/pkg/unifiedrole/export_test.go`, add the role variable
1. In `/services/graph/pkg/unifiedrole/roles_test.go`, extend unit tests to cover the new role
1. In `/services/web/pkg/theme/theme.go`, add the role into the `common.shareRoles`
2. In `/services/graph/README.md`, add the role to the list of built-in roles
