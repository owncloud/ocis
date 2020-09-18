Bugfix: Prevent segfault when no password is set

Passwords are stored in a dedicated child struct of an account. We fixed several segfault conditions where the methods would try to unset a password when that child struct was not existing.

https://github.com/owncloud/ocis-accounts/pull/65
