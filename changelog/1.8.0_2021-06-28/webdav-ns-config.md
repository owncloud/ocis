Bugfix: Make webdav namespace configurable across services

The WebDAV namespace is used across various services, but it was previously
hardcoded in some of the services. This PR uses the same environment variable
to set the config correctly across the services.

https://github.com/owncloud/ocis/pull/2198