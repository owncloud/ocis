Bugfix: Autocreate IDP private key also if file exists but is empty

We've fixed the behavior for the IDP private key generation so that
a private key is also generated when the file already exists but is empty.

https://github.com/owncloud/ocis/pull/4394
