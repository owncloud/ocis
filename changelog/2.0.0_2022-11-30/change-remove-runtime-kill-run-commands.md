Bugfix: Remove runtime kill and run commands

We've removed the kill and run commands from the oCIS runtime.
If these dynamic capabilities are needed, one should switch to a full fledged
supervisor and start oCIS as individual services.

If one wants to start a only a subset of services, this is still possible
by setting OCIS_RUN_EXTENSIONS.

https://github.com/owncloud/ocis/pull/3740
