Bugfix: Do not disable the personal space when the Drives.Create permission check is inconclusive

When a user's role is (re-)assigned, both the proxy (on OIDC login) and the graph
appRoleAssignment handler check the `Drives.Create` permission to decide whether to
restore or disable the user's personal space. The permission check collapsed two very
different outcomes into a single `false`: the user genuinely lacks the permission, and
the permission could not be determined (a transport error, or a non-OK status such as
`CODE_INTERNAL` returned by the settings/gateway service). In the second case the code
proceeded to disable the personal space, moving it to the trash, even though the user's
entitlement was never actually denied.

`utils.CheckPermission` now distinguishes the three cases at the root: a transport error
or a non-OK status other than `PERMISSION_DENIED` is surfaced as an error and the callers
fail closed. The personal space is left untouched, and the graph appRoleAssignment handler
additionally reverts the role assignment it just persisted so the user is not left in a
half-applied state. The role transition is retried on the next login.

https://github.com/owncloud/ocis/issues/12429
