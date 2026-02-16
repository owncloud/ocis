Enhancement: Add logging when a users space gets deleted

When deleting a user, their personal space will also be deleted. When this operation fails the logging in the graph
service was insufficient. We added some logs.

https://github.com/owncloud/ocis/pull/11037
