Bugfix: Fix search service to not log expected cases as errors

We changed the search service to not log cases where resources that were about to be indexed can no longer be found.
Those are expected cases, e.g. when the file in question has already been deleted or renamed meanwhile.

https://github.com/owncloud/ocis/pull/8200
