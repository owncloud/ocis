Enhancement: Limit the CI runners that one pull request uses

The Acceptance Tests workflow now cancels the older run when a new commit
arrives on the same pull request. The API test and e2e test matrices now run
fewer jobs at the same time. One pull request thus cannot use all the GitHub
runners of the organization.

https://github.com/owncloud/ocis/pull/13459
https://github.com/owncloud/ocis/pull/13453
