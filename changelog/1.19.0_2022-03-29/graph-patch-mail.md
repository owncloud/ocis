Bugfix: Fix request validation on GraphAPI User updates

Fix PATCH on graph/v1.0/users when no 'mail' attribute
is present in the request body

https://github.com/owncloud/ocis/issues/3167
