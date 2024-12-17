Bugfix: Bump Reva

Bumps reva version to 2.27.0

*   Bugfix [cs3org/reva#4985](https://github.com/cs3org/reva/pull/4985): Drop unneeded session locks
*   Bugfix [cs3org/reva#5000](https://github.com/cs3org/reva/pull/5000): Fix ceph build
*   Bugfix [cs3org/reva#4989](https://github.com/cs3org/reva/pull/4989): Deleting OCM share also updates storageprovider
*   Enhancement [cs3org/reva#4998](https://github.com/cs3org/reva/pull/4998): Emit event when an ocm share is received
*   Enhancement [cs3org/reva#4996](https://github.com/cs3org/reva/pull/4996): Get rid of some cases of unstructured logging

Bumps reva version to 2.26.8

*   Fix [cs3org/reva#4983](https://github.com/cs3org/reva/pull/4983): Delete stale shares in the jsoncs3 share manager
*   Fix [cs3org/reva#4963](https://github.com/cs3org/reva/pull/4963): Fix name and displayName in an ocm
*   Fix [cs3org/reva#4968](https://github.com/cs3org/reva/pull/4968): Jsoncs3 cache fixes
*   Fix [cs3org/reva#4928](https://github.com/cs3org/reva/pull/4928): Propagate lock in PROPPATCH
*   Fix [cs3org/reva#4971](https://github.com/cs3org/reva/pull/4971): Use manager to list shares
*   Fix [cs3org/reva#4978](https://github.com/cs3org/reva/pull/4978): We added more trace spans in decomposedfs
*   Fix [cs3org/reva#4921](https://github.com/cs3org/reva/pull/4921): Polish propagation related code

Bugfix: Bump Reva to v2.26.7

   * Fix [cs3org/reva#4964](https://github.com/cs3org/reva/issues/4964): Fix a wrong error code when approvider creates a new file

Bump Reva to v2.26.6

   * Fix [cs3org/reva#4955](https://github.com/cs3org/reva/issues/4955): Allow small clock skew in reva token validation
   * Fix [cs3org/reva#4929](https://github.com/cs3org/reva/issues/4929): Fix flaky posixfs integration tests
   * Fix [cs3org/reva#4953](https://github.com/cs3org/reva/issues/4953): Avoid gateway panics
   * Fix [cs3org/reva#4959](https://github.com/cs3org/reva/issues/4959): Fix missing file touched event
   * Fix [cs3org/reva#4933](https://github.com/cs3org/reva/issues/4933): Fix federated sharing when using an external identity provider
   * Fix [cs3org/reva#4935](https://github.com/cs3org/reva/issues/4935): Enable datatx log
   * Fix [cs3org/reva#4936](https://github.com/cs3org/reva/issues/4936): Do not delete mlock files
   * Fix [cs3org/reva#4954](https://github.com/cs3org/reva/issues/4954): Prevent a panic when logging an error
   * Fix [cs3org/reva#4956](https://github.com/cs3org/reva/issues/4956): Improve posixfs error handling and logging
   * Fix [cs3org/reva#4951](https://github.com/cs3org/reva/issues/4951): Pass the initialized logger down the stack

Bugfix: Bump Reva to v2.26.5

   * Fix [cs3org/reva#4926](https://github.com/cs3org/reva/issues/4926): Make etag always match content on downloads
   * Fix [cs3org/reva#4920](https://github.com/cs3org/reva/issues/4920): Return correct status codes for simple uploads
   * Fix [cs3org/reva#4924](https://github.com/cs3org/reva/issues/4924): Fix sync propagation
   * Fix [cs3org/reva#4916](https://github.com/cs3org/reva/issues/4916): Improve posixfs stability and performanc

Enhancement: Bump reva to 2.26.4

*   Bugfix [cs3org/reva#4917](https://github.com/cs3org/reva/pull/4917): Fix 0-byte file uploads
*   Bugfix [cs3org/reva#4918](https://github.com/cs3org/reva/pull/4918): Fix app templates

Bump reva to 2.26.3

*   Bugfix [cs3org/reva#4908](https://github.com/cs3org/reva/pull/4908): Add checksum to OCM storageprovider responses
*   Enhancement [cs3org/reva#4910](https://github.com/cs3org/reva/pull/4910): Bump cs3api
*   Enhancement [cs3org/reva#4909](https://github.com/cs3org/reva/pull/4909): Bump cs3api
*   Enhancement [cs3org/reva#4906](https://github.com/cs3org/reva/pull/4906): Bump cs3api


Bump reva to 2.26.2

*   Enhancement [cs3org/reva#4897](https://github.com/cs3org/reva/pull/4897): Fix remaining quota calculation
*   Bugfix      [cs3org/reva#4902](https://github.com/cs3org/reva/pull/4902): Fix quota calculation

https://github.com/owncloud/ocis/pull/10766
https://github.com/owncloud/ocis/pull/10735
https://github.com/owncloud/ocis/pull/10612
https://github.com/owncloud/ocis/pull/10552
https://github.com/owncloud/ocis/pull/10539
https://github.com/owncloud/ocis/pull/10419
