Enhancement: Add vault capabilities to the OCS capabilities endpoint

Added `OCIS_ENABLE_VAULT_MODE` / `FRONTEND_ENABLE_VAULT_MODE` config option to
the frontend service. When enabled, the OCS capabilities endpoint advertises
`vault.enabled = true`. Clients can request vault-specific capabilities via
`/ocs/v2.php/cloud/capabilities?vault=true`, which returns a response with
public sharing and federation sharing disabled.

https://github.com/owncloud/ocis/pull/12283
https://github.com/owncloud/reva/pull/584
