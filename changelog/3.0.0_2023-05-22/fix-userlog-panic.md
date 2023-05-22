Bugfix: Fix userlog panic

userlog services paniced because of `nil` ctx. That is fixed now

https://github.com/owncloud/ocis/pull/6114
