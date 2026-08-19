Security: Make OIDC group sync opt-in

**What changed.** `PROXY_AUTOPROVISION_CLAIM_GROUPS` now defaults to `""`, which disables OIDC group membership sync (and, with it, creation of local groups from claim values). It previously defaulted to `groups`. Setting it to a non-empty claim name restores the previous behaviour unchanged: memberships are synced and groups named in the claim are created if they do not exist locally.

**Why.** With `PROXY_AUTOPROVISION_ACCOUNTS=true` and the previous `groups` default, the proxy synced group memberships from the OIDC `groups` claim on every authenticated request out of the box, creating local groups for any claim value that did not already exist. In identity providers that let ordinary users create groups with arbitrary names, this allowed an unprivileged user to inject group names into oCIS. Defaulting the claim to empty makes group sync an explicit opt-in.

**Upgrade note.** Deployments that set `PROXY_AUTOPROVISION_CLAIM_GROUPS` explicitly are unaffected. Deployments that relied on the previous `groups` default without setting it will stop syncing group memberships after upgrade; set `PROXY_AUTOPROVISION_CLAIM_GROUPS=groups` to restore the previous behaviour.

Note: matching claim values to existing local groups is still done by display name. Hardening that matching is tracked separately and is not part of this change.

https://github.com/owncloud/ocis/pull/12490
