Enhancement: Modify the concurrency default

We have changed the default MaxConcurrency value from 100 to 5 to prevent too frequent gc runs on low memory systems.
We have also bumped reva to pull in the related changes from there.

https://github.com/owncloud/ocis/pull/8309
https://github.com/owncloud/ocis/issues/8257
https://github.com/cs3org/reva/pull/4485
