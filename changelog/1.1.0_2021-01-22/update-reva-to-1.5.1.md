Enhancement: Update reva to v1.5.1

Summary
-------

* Fix #1401: Use the user in request for deciding the layout for non-home DAV requests
* Fix #1413: Re-include the '.git' dir in the Docker images to pass the version tag
* Fix #1399: Fix ocis trash-bin purge
* Enh #1397: Bump the Copyright date to 2021
* Enh #1398: Support site authorization status in Mentix
* Enh #1393: Allow setting favorites, mtime and a temporary etag
* Enh #1403: Support remote cloud gathering metrics

Details
-------

* Bugfix #1401: Use the user in request for deciding the layout for non-home DAV requests

  For the incoming /dav/files/userID requests, we have different namespaces depending on
  whether the request is for the logged-in user's namespace or not. Since in the storage drivers,
  we specify the layout depending only on the user whose resources are to be accessed, this fails
  when a user wants to access another user's namespace when the storage provider depends on the
  logged in user's namespace. This PR fixes that.

  For example, consider the following case. The owncloud fs uses a layout {{substr 0 1
  .Id.OpaqueId}}/{{.Id.OpaqueId}}. The user einstein sends a request to access a resource
  shared with him, say /dav/files/marie/abcd, which should be allowed. However, based on the
  way we applied the layout, there's no way in which this can be translated to /m/marie/.

  https://github.com/cs3org/reva/pull/1401

* Bugfix #1413: Re-include the '.git' dir in the Docker images to pass the version tag

  And git SHA to the release tool.

  https://github.com/cs3org/reva/pull/1413

* Bugfix #1399: Fix ocis trash-bin purge

  Fixes the empty trash-bin functionality for ocis-storage

  https://github.com/owncloud/product/issues/254
  https://github.com/cs3org/reva/pull/1399

* Enhancement #1397: Bump the Copyright date to 2021

  https://github.com/cs3org/reva/pull/1397

* Enhancement #1398: Support site authorization status in Mentix

  This enhancement adds support for a site authorization status to Mentix. This way, sites
  registered via a web app can now be excluded until authorized manually by an administrator.

  Furthermore, Mentix now sets the scheme for Prometheus targets. This allows us to also support
  monitoring of sites that do not support the default HTTPS scheme.

  https://github.com/cs3org/reva/pull/1398

* Enhancement #1393: Allow setting favorites, mtime and a temporary etag

  We now let the oCIS driver persist favorites, set temporary etags and the mtime as arbitrary
  metadata.

  https://github.com/owncloud/ocis/issues/567
  https://github.com/cs3org/reva/issues/1394
  https://github.com/cs3org/reva/pull/1393

* Enhancement #1403: Support remote cloud gathering metrics

  The current metrics package can only gather metrics either from json files. With this feature,
  the metrics can be gathered polling the http endpoints exposed by the owncloud/nextcloud
  sciencemesh apps.

  https://github.com/cs3org/reva/pull/1403

https://github.com/owncloud/ocis/pull/1372
