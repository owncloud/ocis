Bugfix: Adjust UUID validation to be more tolerant

The UUID now allows any alphanumeric character and "-", "_", ".", "+" and "@" which
can also allow regular user names.

https://github.com/owncloud/ocis-settings/issues/41
