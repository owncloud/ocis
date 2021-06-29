Enhancement: Properly configure graph-explorer client registration

The client registration in the `identifier-registration.yaml` for the graph-explorer didn't contain `redirect_uris` nor `origins`. Both were added to prevent exploitation.

https://github.com/owncloud/ocis/pull/2118
