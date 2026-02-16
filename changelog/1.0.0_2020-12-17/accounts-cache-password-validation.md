Change: Cache password validation

Tags: accounts

The password validity check for requests like `login eq '%s' and password eq '%s'` is now cached for 10 minutes.
This improves the performance for basic auth requests. 

https://github.com/owncloud/ocis/pull/958