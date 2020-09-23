Change: We make use of the roles manager to enforce permission checks

The roles cache and its cache update middleware have been replaced with a roles manager in ocis-pkg/v2. We've switched
over to the new roles manager implementation, to prepare for permission checks on grpc requests as well.

<https://github.com/owncloud/ocis/accounts/pull/108>
<https://github.com/owncloud/ocis-pkg/pull/60>
