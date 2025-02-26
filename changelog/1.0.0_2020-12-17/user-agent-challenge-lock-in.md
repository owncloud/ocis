Enhancement: Add www-authenticate based on user agent

Tags: reva, proxy

We now comply with HTTP spec by adding Www-Authenticate headers on every `401` request. Furthermore, we not only take care of such a thing at the Proxy but also Reva will take care of it. In addition, we now are able to lock-in a set of User-Agent to specific challenges.

Admins can use this feature by configuring oCIS + Reva following this approach:

```
STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT="mirall:basic, Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:83.0) Gecko/20100101 Firefox/83.0:bearer" \
PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT="mirall:basic, Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:83.0) Gecko/20100101 Firefox/83.0:bearer" \
PROXY_ENABLE_BASIC_AUTH=true \
go run cmd/ocis/main.go server
```

We introduced two new environment variables:

`STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT` as well as `PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT`, The reason they have the same value is not to rely on the os env on a distributed environment, so in redundancy we trust. They both configure the same on the backend storage and oCIS Proxy.

https://github.com/owncloud/ocis/pull/1009
