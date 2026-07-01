Enhancement: Make OIDC group sync opt-in and disable claim-driven group creation

**What changed.**

- `PROXY_AUTOPROVISION_CLAIM_GROUPS` now defaults to `""`, which disables OIDC group membership sync. It previously defaulted to `groups`.
- The new `PROXY_AUTOPROVISION_GROUP_CREATE` flag (default `false`) controls whether a local group is created when the groups claim contains a name that does not exist locally. When disabled, such a claim value is skipped instead of creating a group.

**Why.** With `PROXY_AUTOPROVISION_ACCOUNTS=true` and the previous `groups` default, the proxy synced group memberships from the OIDC `groups` claim on every authenticated request out of the box, creating local groups for any claim value that did not already exist. In identity providers that let ordinary users create groups with arbitrary names, this allowed an unprivileged user to inject group names into oCIS. Making sync opt-in and not creating groups from claims by default removes this exposure under the default configuration.

**Upgrade note.** Deployments that relied on the previous `groups` default and did not set `PROXY_AUTOPROVISION_CLAIM_GROUPS` explicitly will stop syncing group memberships after upgrade. To keep syncing, set `PROXY_AUTOPROVISION_CLAIM_GROUPS=groups`; to also keep creating groups from claim values as before, additionally set `PROXY_AUTOPROVISION_GROUP_CREATE=true`.

Note: matching claim values to existing local groups is still done by display name. Hardening that matching is tracked separately and is not part of this change.

https://github.com/owncloud/ocis/pull/12490
