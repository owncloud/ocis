Enhancement: Add numeric uid and gid to the access token

The eos storage driver is fetching the uid and gid of a user from the access token. This PR is using the response of the accounts service to mint them in the token.

https://github.com/owncloud/ocis-proxy/pull/89
