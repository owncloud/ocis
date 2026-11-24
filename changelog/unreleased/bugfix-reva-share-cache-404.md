Bugfix: Remove extra storage calls for missing share cache files

The share manager does not look at the trash of the share metadata space
anymore when a share cache file is missing. This lookup could never find the
file, and it cost two more storage calls for each missing file. This applies to
the created shares, to the received shares of users and groups, and to the
shares of a space.

The share manager now creates an empty share cache file only in a listing. Add
writes the file with the share in one upload. A removal does not create the
file when it removes nothing.

https://github.com/owncloud/ocis/pull/13565
https://github.com/owncloud/reva/pull/845
https://github.com/owncloud/reva/pull/848
