Security: Bump crewjam/saml and gobwas/ws dependencies

We've bumped github.com/crewjam/saml to v0.5.2 to fix GO-2026-6909, a nil
pointer panic in validateAssertion, and github.com/gobwas/ws to v1.5.0 to fix
GO-2026-6940, an unbounded allocation caused by a peer-declared frame length.

https://github.com/owncloud/ocis/pull/13558
