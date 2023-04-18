Enhancement: Add optional services to the runtime

Make it possible to start optional services in the ocis runtime. Instead of using `OCIS_RUN_SERVICES` to define all services we can now use `OCIS_ADD_RUN_SERVICES` to add a comma separated list of additional services which are not started in the single process by default.

https://github.com/owncloud/ocis/pull/6071
