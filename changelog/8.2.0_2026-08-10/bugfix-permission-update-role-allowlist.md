Bugfix: Apply the role allowlist on permission updates

The `PATCH` permission and space-root permission handlers now consult the
administrator role allowlist when validating a request, consistent with the
invite handler.

https://github.com/owncloud/ocis/pull/12540
