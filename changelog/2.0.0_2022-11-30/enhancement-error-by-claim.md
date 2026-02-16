Enhancement: Improve error log for "could not get user by claim" error

We've improved the error log for "could not get user by claim" error where
previously only the "nil" error has been logged. Now we're logging the
message from the transport.

https://github.com/owncloud/ocis/pull/4227
