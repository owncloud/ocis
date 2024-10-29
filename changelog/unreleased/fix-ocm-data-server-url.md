Bugfix: Allow to configure data server URL for ocm

We introduced the `OCM_OCM_STORAGE_DATA_SERVER_URL` setting to fix a bug
when downloading files from an OCM share. Before the data server URL defaulted
to the listen address of the OCM server, which did not work when using
0.0.0.0 as the listen address.

https://github.com/owncloud/ocis/pull/10440
https://github.com/owncloud/ocis/issues/10358

