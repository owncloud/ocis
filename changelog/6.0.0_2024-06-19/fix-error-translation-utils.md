Bugfix: Fix the error translation from utils

We've fixed the error translation from the statusCodeError type to CS3 Status because the FromCS3Status function converts a CS3 status code into a corresponding local Error representation.

https://github.com/owncloud/ocis/pull/9331
https://github.com/owncloud/ocis/issues/9151
