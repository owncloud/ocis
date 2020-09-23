Change: Default to config based user backend

We changed the default configuration to use the config file backend instead of the ownCloud backend.

The config backend currently only has two hard coded users: demo and admin. To switch back to the ownCloud backend use `GLAUTH_BACKEND_DATASTORE=owncloud`

<https://github.com/owncloud/ocis/glauth/pull/6>
