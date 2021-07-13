Enhancement: update REVA to v1.9.1

* Fix cs3org/reva#1843: Correct Dockerfile path for the reva CLI and alpine3.13 as builder
* Fix cs3org/reva#1835: Cleanup owncloudsql driver
* Fix cs3org/reva#1868: Minor fixes to the grpc/http plugin: checksum, url escaping
* Fix cs3org/reva#1885: Fix template in eoshomewrapper to use context user rather than resource
* Fix cs3org/reva#1833: Properly handle name collisions for deletes in the owncloud driver
* Fix cs3org/reva#1874: Use the original file mtime during upload
* Fix cs3org/reva#1854: Add the uid/gid to the url for eos
* Fix cs3org/reva#1848: Fill in missing gid/uid number with nobody
* Fix cs3org/reva#1831: Make the ocm-provider endpoint in the ocmd service unprotected
* Fix cs3org/reva#1808: Use empty array in OCS Notifications endpoints
* Fix cs3org/reva#1825: Raise max grpc message size
* Fix cs3org/reva#1828: Send a proper XML header with error messages
* Chg cs3org/reva#1828: Remove the oidc provider in order to upgrad mattn/go-sqlite3 to v1.14.7
* Enh cs3org/reva#1834: Add API key to Mentix GOCDB connector
* Enh cs3org/reva#1855: Minor optimization in parsing EOS ACLs
* Enh cs3org/reva#1873: Update the EOS image tag to be for revad-eos image
* Enh cs3org/reva#1802: Introduce list spaces
* Enh cs3org/reva#1849: Add readonly interceptor
* Enh cs3org/reva#1875: Simplify resource comparison
* Enh cs3org/reva#1827: Support trashbin sub paths in the recycle API

https://github.com/owncloud/ocis/pull/2280
