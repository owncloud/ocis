Bugfix: Fix the access to graph users and drive endpoints

Restrict the access to `GET /graph/v1.0/users/{userID}` and `GET /graph/v1.0/users/{userID}/drive`
endpoints for the user without the admin permission.

https://github.com/owncloud/ocis/issues/12347
