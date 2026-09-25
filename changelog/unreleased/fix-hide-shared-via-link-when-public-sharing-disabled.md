Bugfix: Hide "Shared via link" when public sharing is disabled

The Shares navigation (the desktop tab list and the mobile drop-down) always
offered a "Shared via link" entry, even when public link sharing was disabled
for the account — for example in the Vault, where the backend rejects public
link creation and the capability `files_sharing.public.enabled` reported
`false`. The entry led to a page that could never contain anything.

The Shares navigation was changed to hide the "Shared via link" entry
whenever `files_sharing.public.enabled` is `false`. Anyone who still reached
the via-link route directly, including through the legacy
`/list/shared-via-link` redirect, was redirected to "Shared with me" at
navigation time, keeping the current route scope (e.g. `vault`) intact. When
public sharing was enabled, all three Shares navigation entries continued to
behave exactly as before.

**PR:** [#13005](https://github.com/owncloud/ocis/pull/13005)
