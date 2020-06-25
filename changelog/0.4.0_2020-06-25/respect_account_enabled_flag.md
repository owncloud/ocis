Enhancement: respect account_enabled flag

If the account returned by the accounts service has the account_enabled flag
set to false, the proxy will return immediately with the status code unauthorized.

https://github.com/owncloud/ocis-proxy/issues/53
