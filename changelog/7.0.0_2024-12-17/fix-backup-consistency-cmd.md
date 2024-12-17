Bugfix: 'ocis backup consistency' fixed for file revisions

A bug was fixed that caused the 'ocis backup consistency' command to incorrectly report
inconistencies when file revisions with a zero value for the nano-second part of the
timestamp were present.

https://github.com/owncloud/ocis/pull/10493
https://github.com/owncloud/ocis/issues/9498
