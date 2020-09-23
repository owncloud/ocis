Enhancement: create account if it doesn't exist in ocis-accounts

The accounts_uuid middleware tries to get the account from ocis-accounts.
If it doens't exist there yet the proxy creates the account using the ocis-account api.

<https://github.com/owncloud/ocis/proxy/issues/55>
<https://github.com/owncloud/ocis/proxy/issues/58>
