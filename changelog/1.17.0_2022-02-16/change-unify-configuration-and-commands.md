Change: Unify configuration and commands

We've unified the configuration and commands of all non storage services. This also
includes the change, that environment variables are now defined on the config struct
as tags instead in a separate mapping.

https://github.com/owncloud/ocis/pull/2818
