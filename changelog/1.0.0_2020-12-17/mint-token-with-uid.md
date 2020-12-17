Bugfix: Mint token with uid and gid

Tags: accounts

The eos driver expects the uid and gid from the opaque map of a user. While the proxy does mint tokens correctly, the accounts service wasn't.

https://github.com/owncloud/ocis/pull/737

