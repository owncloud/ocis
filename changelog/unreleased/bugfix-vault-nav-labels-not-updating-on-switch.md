Bugfix: Vault navigation labels not updating on switch

The Files app navigation labels for the personal space and spaces list
sometimes kept showing the regular names instead of the vault-specific
ones after switching into the vault, and vice versa. This happened
because the vault detection only looked at the plain URL path and
missed cases such as hash-based routes or the redirect right after
logging in.

The vault detection has been centralized and now also recognizes
hash-based vault routes and the post-login redirect target, so the
navigation labels correctly reflect whether the user is in the vault.

https://github.com/owncloud/ocis/pull/12729
