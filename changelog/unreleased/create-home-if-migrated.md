Enhancement: only send create home request if an account has been migrated 

This change adds a check if an account has been migrated by getting it from the
ocis-accounts service. If no account is returned it means it hasn't been migrated.

https://github.com/owncloud/ocis-proxy/issues/52
