# Webfinger service

The webfinger service provides an RFC7033 WebFinger lookup of ownCloud instances relevant for a given user account.

It is based on https://github.com/owncloud/lookup-webfinger-sciebo but also returns a `displayname` in addition to the `href` property.

# Links
Initial issue: https://github.com/owncloud/ocis/issues/5281

# TODO
- [ ] actually query a backend, for now ldap in context of a multi tenant deployment ... or use the graph api to list /education/schools for a user? ldap would be more direct.
- [ ] don't use the broken metrics/instrumentation/logging interface from graph / webdav services, see https://github.com/owncloud/ocis/issues/5209

