Bugfix: return proper errors when ocs/cloud/users is using the cs3 backend

The ocs API was just exiting with a fatal error on any update request,
when configured for the cs3 backend. Now it returns a proper error.

https://github.com/owncloud/ocis/issues/3483
