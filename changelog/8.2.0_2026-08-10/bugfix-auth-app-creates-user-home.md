Bugfix: The auth-app will create the user's home if needed

When the user logs in, his home must be created. This happens automatically
during login via OIDC (web access). Some recent changes in the code broke
this behavior when the user logs in via auth-app.
Now, this behavior is restored, and the user's home will be created when the
user logs in via auth-app.

https://github.com/owncloud/ocis/pull/12457
