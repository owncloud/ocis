Enhancement: Add option to disable direct (user/group) sharing

Added an `OCIS_ENABLE_USER_SHARING` config option, read by the frontend,
graph, sharing and ocm services. It defaults to `true`. When set to
`false`, creating new user, group or federated shares is rejected, the
legacy `sharees` search endpoint returns no results, and the
`files_sharing.user.enabled` and `files_sharing.user_enumeration.enabled`
capabilities report `false` so clients hide the corresponding UI. Public
link sharing and space membership are not affected.

https://github.com/owncloud/ocis/pull/12542
