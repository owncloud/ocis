Enhancement: Force Strict-Transport-Security

Added `PROXY_FORCE_STRICT_TRANSPORT_SECURITY` environment variable to force emission of `Strict-Transport-Security` header on all responses, including plain HTTP requests when TLS is terminated upstream. Useful when oCIS is deployed behind a proxy.

https://github.com/owncloud/ocis/pull/11880
