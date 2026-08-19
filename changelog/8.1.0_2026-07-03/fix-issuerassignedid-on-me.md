Bugfix: Return correct issuerAssignedId on /me

The `/graph/v1.0/me` endpoint reported the internal user UUID as
`identities[].issuerAssignedId` instead of the issuer-assigned identity (the
OIDC `sub`). The endpoint took a fast path that built the user model from the
CS3 user in the request context, which does not carry the external identity, so
it fell back to the internal UUID. `/me` now always resolves the user through
the identity backend, which reads the stored external identity and returns the
correct value. Group memberships are still only expanded when
`$expand=memberOf` is requested.

https://github.com/owncloud/ocis/pull/12411
