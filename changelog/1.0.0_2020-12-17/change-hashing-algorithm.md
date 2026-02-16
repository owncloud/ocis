Change: Use bcrypt to hash the user passwords 

Change the hashing algorithm from SHA-512 to bcrypt since the latter is better suitable for password hashing.
This is a breaking change. Existing deployments need to regenerate the accounts folder.


https://github.com/owncloud/ocis/issues/510
