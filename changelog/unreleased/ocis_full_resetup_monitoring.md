Bugfix: Rework monitoring in the ocis_full deployment example

The ocis_full deployment example has been basically reworked for how to provide monitoring.

We now have:
- a singe place for the definition of the tracing envvars for all ocis related container services
- an easy and modular setup defining which sources should be inlcuded in monitoring via .env
- comments describing the setup for the ease extending it

https://github.com/owncloud/ocis/pull/11995
