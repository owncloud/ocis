Enhancement: Allow scaling the search service

Previously, the search service locked the index for its whole lifetime,
so any other search service wouldn't be able to access to the index. With this
change, the search service can be configure to lock the index per operation,
so other search services can access the index as long as there is no operation
ongoing.

https://github.com/owncloud/ocis/pull/11029
