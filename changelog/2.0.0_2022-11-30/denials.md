Enhancement: Deny access to resources

We added an experimental feature to deny access to a certain resource. This feature is disabled by default and considered as EXPERIMENTAL. You can enable it by setting FRONTEND_OCS_ENABLE_DENIALS to `true`. It announces an available deny access permission via WebDAV on each resource. By convention it is only possible to deny access on folders. The clients can check the presence of the feature by the capability `deny_access` in the `files_sharing` section.

https://github.com/owncloud/ocis/pull/4903
