Change: CS3 can be used as accounts-backend

Tags: proxy

PROXY_ACCOUNT_BACKEND_TYPE=cs3
PROXY_ACCOUNT_BACKEND_TYPE=accounts (default)

By using a backend which implements the CS3 user-api (currently provided by reva/storage) it is possible to bypass
the ocis-accounts service and for example use ldap directly.


https://github.com/owncloud/ocis/pull/1020
