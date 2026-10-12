Bugfix: Allow space managers to remove shares created by other users

Space managers with permission to manage a resource's shares were unable to
remove shares that other users had created on that resource. Attempting to
remove such a share failed with a generic server error instead of succeeding.

Space managers can now remove any share on a resource within their space,
regardless of who created it, as long as they hold the required permission
on that resource.

https://github.com/owncloud/ocis/pull/12705
