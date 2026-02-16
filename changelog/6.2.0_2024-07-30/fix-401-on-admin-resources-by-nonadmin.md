Bugfix: Fix restarting of postprocessing

We fixed a bug where non-admin requests to admin resources would get 401 Unauthorized.
Now, the server sends 403 Forbidden response.

https://github.com/owncloud/ocis/pull/6945
https://github.com/owncloud/ocis/issues/5938
