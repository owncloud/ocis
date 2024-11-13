Bugfix: Set MaxConcurrency to 1

Set MaxConcurrency for frontend and userlog services to 1. Too many workers will negatively impact performance on small machines.

https://github.com/owncloud/ocis/pull/10557
