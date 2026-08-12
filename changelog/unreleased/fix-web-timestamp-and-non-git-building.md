Bugfix: Generate web timestamp by env/git log; make it buildable without .git

Change compilationTimeStamp to distTimeStamp, set its value from env SOURCE_DATE_EPOCH, or timestamp of last git commit, or fallback to current time. Also fixes a problem that raises an error during make ci-node-generate due to running git clean -xfd in a non-repository directory (usually an extracted source tarball), showing a message instead.

https://github.com/owncloud/ocis/pull/12758
