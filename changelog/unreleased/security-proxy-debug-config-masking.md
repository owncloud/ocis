Security: Mask secret values in the proxy debug output

The debug server of the proxy service did not mask all secret configuration
values. The proxy now masks these values. A new test checks that every secret
value of the proxy configuration stays masked.

https://github.com/owncloud/ocis/pull/13220
