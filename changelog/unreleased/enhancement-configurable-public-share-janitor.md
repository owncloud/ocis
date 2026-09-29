Enhancement: Make the public share expiry janitor configurable

The sharing service runs a background janitor that permanently deletes
expired public shares. Whether that cleanup runs at all could previously only
be set in the sharing service's yaml config, and how often it ran was fixed
internally with no way to tune it.

Both settings are now exposed as environment variables:
OCIS_SHARING_ENABLE_EXPIRED_SHARES_CLEANUP toggles the cleanup (default:
enabled, expired shares stay hidden from listings even when disabled), and
the new OCIS_SHARING_JANITOR_RUN_INTERVAL sets the interval in seconds
between janitor runs (default: 600).

https://github.com/owncloud/ocis/pull/13707
