Enhancement: Unified Roles Management

Improved management of unified roles with the introduction of default enabled/disabled states and a new command for listing available roles.
It is important to note that a disabled role does not lose previously assigned permissions;
it only means that the role is not available for new assignments.

The following roles are now enabled by default:

- UnifiedRoleViewerID
- UnifiedRoleSpaceViewer
- UnifiedRoleEditor
- UnifiedRoleSpaceEditor
- UnifiedRoleFileEditor
- UnifiedRoleEditorLite
- UnifiedRoleManager

The following roles are now disabled by default:

- UnifiedRoleSecureViewer

To enable the UnifiedRoleSecureViewer role, you must provide a list of all available roles through one of the following methods:

- Using the GRAPH_AVAILABLE_ROLES environment variable.
- Setting the available_roles configuration value.

To enable a role, include the UID of the role in the list of available roles.

A new command has been introduced to simplify the process of finding out which UID belongs to which role. The command is:

```
$ ocis graph list-unified-roles
```

The output of this command includes the following information for each role:

- uid: The unique identifier of the role.
- Description: A short description of the role.
- Enabled: Whether the role is enabled or not.

https://github.com/owncloud/ocis/pull/9727
https://github.com/owncloud/ocis/issues/9698
