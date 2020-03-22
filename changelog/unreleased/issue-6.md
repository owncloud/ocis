Change: start multiple services with dedicated commands

The initial version would only allow us to use a set of reva configurations to start multiple services.
We use a more opinionated set of commands to start dedicated services that allows us to configure them individually.
It allows us to switch eg. the user backend to LDAP and fully use it on the cli.

https://github.com/owncloud/ocis-reva/issues/6