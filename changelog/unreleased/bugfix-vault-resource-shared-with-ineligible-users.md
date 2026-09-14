Bugfix: Do not share vault resources with users who cannot access them

A vault ("Safe") resource could be shared with a user lacking the "Vault mode"
permission. Such a user can never enter the vault, so the share was useless to
them, while its name was still disclosed to them.

Sharing a vault resource now requires the grantee to hold that permission, and
the check fails closed. The recipient picker no longer offers ineligible users,
nor groups at all, since the permission is assigned per account.

https://github.com/owncloud/ocis/pull/12905
