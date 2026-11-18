Bugfix: Disable async uploads in litmus CI to prevent 425 Too Early errors

We disabled async uploads in the litmus CI runner because they caused
intermittent 425 Too Early failures.

https://github.com/owncloud/ocis/pull/13519
