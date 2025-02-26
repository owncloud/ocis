Enhancement: Allow providing list of services NOT to start

Until now if one wanted to use a custom version of a service, one
needed to provide `OCIS_RUN_SERVICES` which is a list of all services to start.
Now one can provide `OCIS_EXCLUDE_RUN_SERVICES` which is a list of only services not to start

https://github.com/owncloud/ocis/pull/4254
