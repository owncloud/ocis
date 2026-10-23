Bugfix: Disable confirm in the link password dialog until a password is entered

The "Add password" / "Edit password" dialog for links could be confirmed with
an empty field, which removed the password from the link. The confirm button
is now disabled until a valid password is entered.

https://github.com/owncloud/ocis/pull/13218
