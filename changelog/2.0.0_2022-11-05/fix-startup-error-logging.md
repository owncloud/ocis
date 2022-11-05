Bugfix: Fix startup error logging

We've fixed the startup error logging, so that users will the reason for a failed
startup even on "error" log level. Previously they would only see it on "info" log level.
Also in a lot of cases the reason for the failed shutdown was omitted.

https://github.com/owncloud/ocis/pull/4093
