Bugfix: Speed up the admin Users page

The admin Users page could take tens of seconds to load because roles were
looked up one user at a time. Role assignments for a full page of users are
now fetched in a single batched, concurrent request.

https://github.com/owncloud/ocis/pull/13501
