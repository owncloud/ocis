Bugfix: Ignore write error after successful read

In the settings service, we're caching data in the NATSjs store. When reading
service information, that information is written in the cache. Previously,
failing to write that information in the cache caused an error, despite
correct information being retrieved from the service successfully. Now, if the
information isn't written in the cache, we log the failure, but return the
read information without error.

https://github.com/owncloud/ocis/pull/12772
