Enhancement: Default 30-day expiration for vault shares

When a user or group share of a vault file or folder is created without an
expiration date, the server now preconfigures a default expiration (30 days
by default, configurable via `GRAPH_DEFAULT_VAULT_SHARE_EXPIRATION_DAYS`) to
encourage data minimization. The sharer can override or shorten this value.

Non-vault shares and space memberships are exempt and stay unbounded, as
space memberships represent a permanent organizational role.

https://github.com/owncloud/ocis/pull/12988
