Bugfix: Keep group memberships when the OIDC groups claim is absent

When auto-provisioning group memberships from an OIDC claim, the proxy
reconciled the user's groups against the groups claim and removed them from any
group not present in the claim. If a token carried no groups claim at all — for
example a token issued for an OIDC client that has no groups mapper configured,
such as a dedicated desktop-client registration — the parsed group set was empty
and the user was removed from all of their groups.

The sync is now skipped when the groups claim is absent or null in the token. A
present-but-empty claim is still treated as a legitimate "no groups" and
reconciled as before. This mirrors the guard the role-assignment path already
has for a missing roles claim.

https://github.com/owncloud/ocis/issues/11435
https://github.com/owncloud/ocis/pull/12420
