Bugfix: Return an error if we can't get the keys and ensure they're cached

Previously, there was an issue where we could get an error while getting the
public keys from the /hosting/discovery endpoint but we're returning a wrong
success value instead. This is fixed now and we're returning the error.

In addition, the public keys weren't being cached, so we hit the
/hosting/discovery endpoint every time we need to use the public keys. The keys
are now cached so we don't need to hit the endpoint more than what we need.

https://github.com/owncloud/ocis/pull/10590
