Enhancement: Allow to resolve public shares without the ocs tokeninfo endpoint


Instead of querying the /v1.php/apps/files_sharing/api/v1/tokeninfo/ endpoint, a client can now resolve public and internal links by sending a PROPFIND request to /dav/public-files/{sharetoken}

* authenticated clients accessing an internal link are redirected to the "real" resource (`/dav/spaces/{target-resource-id}
* authenticated clients are able to resolve public links like before. For password protected links they need to supply the password even if they have access to the underlying resource by other means.
* unauthenticated clients accessing an internal link get a 401 returned with  WWW-Authenticate set to Bearer (so that the client knows that it need to get a token via the IDP login page.
* unauthenticated clients accessing a password protected link get a 401 returned with an error message to indicate the requirement for needing the link's password.

https://github.com/owncloud/ocis/pull/8926
https://github.com/cs3org/reva/pull/4653
https://github.com/owncloud/ocis/issues/8858
