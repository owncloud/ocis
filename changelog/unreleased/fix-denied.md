Bugfix: Fix deny access for graph roles

We added a unified role "Cannot access" to prevent a regression when switching the share implementation to the graph API. This role is now used to deny access to a resource.The new role is not enabled by default. The whole deny feature is still experimental.

https://github.com/owncloud/ocis/pull/10627
