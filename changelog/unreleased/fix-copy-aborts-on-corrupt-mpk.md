Bugfix: Skip corrupt 0-byte .mpk metadata files instead of aborting list operations

Listing a folder that contained a node with a 0-byte `.mpk` metadata file
previously aborted the entire list operation, which could break a COPY of a
large folder structure partway through. The metadata read now returns an
error for empty `.mpk` files, allowing the existing skip-on-error logic in
the tree walker to log a warning and skip the broken node instead of
aborting the listing.

https://github.com/owncloud/reva/pull/599
https://github.com/owncloud/ocis/pull/12326
