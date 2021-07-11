Enhancement: update REVA to v1.9.1

* Fix #1843: Correct Dockerfile path for the reva CLI and alpine3.13 as builder
* Fix #1835: Cleanup owncloudsql driver
* Fix #1868: Minor fixes to the grpc/http plugin: checksum, url escaping
* Fix #1885: Fix template in eoshomewrapper to use context user rather than resource
* Fix #1833: Properly handle name collisions for deletes in the owncloud driver
* Fix #1874: Use the original file mtime during upload
* Fix #1854: Add the uid/gid to the url for eos
* Fix #1848: Fill in missing gid/uid number with nobody
* Fix #1831: Make the ocm-provider endpoint in the ocmd service unprotected
* Fix #1808: Use empty array in OCS Notifications endpoints
* Fix #1825: Raise max grpc message size
* Fix #1828: Send a proper XML header with error messages
* Chg #1828: Remove the oidc provider in order to upgrad mattn/go-sqlite3 to v1.14.7
* Enh #1834: Add API key to Mentix GOCDB connector
* Enh #1855: Minor optimization in parsing EOS ACLs
* Enh #1873: Update the EOS image tag to be for revad-eos image
* Enh #1802: Introduce list spaces
* Enh #1849: Add readonly interceptor
* Enh #1875: Simplify resource comparison
* Enh #1827: Support trashbin sub paths in the recycle API

https://github.com/owncloud/ocis/pull/2280
