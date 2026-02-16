Bugfix: use same jwt secret for accounts as for metadata storage

We've the metadata storage uses the same jwt secret as all other REVA services.
Therefore the accounts service needs to use the same secret.

Secrets are documented here: https://owncloud.dev/ocis/deployment/#change-default-secrets

https://github.com/owncloud/ocis/pull/3081
