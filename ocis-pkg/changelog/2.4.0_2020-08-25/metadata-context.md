Change: use go-micro's metadata context for account id

We switched to using go-micro's metadata context for reliably passing the AccountID in the context
across service boundaries.

<https://github.com/owncloud/ocis-pkg/pull/56>
