Enhancement: update REVA

Update REVA from v1.10.1-0.20210730095301-fcb7a30a44a6 to v1.11.1-0.20210809134415-3fe79c870fb5
* Fix cs3org/reva#1978: Fix owner type is optional
* Fix cs3org/reva#1965: fix value of file_target in shares
* Fix cs3org/reva#1960: fix updating shares in the memory share manager
* Fix cs3org/reva#1956: fix trashbin listing with depth 0
* Fix cs3org/reva#1957: fix etag propagation on deletes
* Enh cs3org/reva#1861: [WIP] Runtime plugins
* Fix cs3org/reva#1954: fix response format of the sharees API
* Fix cs3org/reva#1819: Remove notifications key from ocs response
* Enh cs3org/reva#1946: Add a share manager that connects to oc10 databases
* Fix cs3org/reva#1899: Fix chunked uploads for new versions
* Fix cs3org/reva#1906: Fix copy over existing resource
* Fix cs3org/reva#1891: Delete Shared Resources as Receiver
* Fix cs3org/reva#1907: Error when creating folder with existing name
* Fix cs3org/reva#1937: Do not overwrite more specific matches when finding storage providers
* Fix cs3org/reva#1939: Fix the share jail permissions in the decomposedfs
* Fix cs3org/reva#1932: Numerous fixes to the owncloudsql storage driver
* Fix cs3org/reva#1912: Fix response when listing versions of another user
* Fix cs3org/reva#1910: Get user groups recursively in the cbox rest user driver
* Fix cs3org/reva#1904: Set Content-Length to 0 when swallowing body in the datagateway
* Fix cs3org/reva#1911: Fix version order in propfind responses
* Fix cs3org/reva#1926: Trash Bin in oCIS Storage Operations
* Fix cs3org/reva#1901: Fix response code when folder doesnt exist on upload
* Enh cs3org/reva#1785: Extend app registry with AddProvider method and mimetype filters
* Enh cs3org/reva#1938: Add methods to get and put context values
* Enh cs3org/reva#1798: Add support for a deny-all permission on references
* Enh cs3org/reva#1916: Generate updated protobuf bindings for EOS GRPC
* Enh cs3org/reva#1887: Add "a" and "l" filter for grappa queries
* Enh cs3org/reva#1919: Run gofmt before building
* Enh cs3org/reva#1927: Implement RollbackToVersion for eosgrpc (needs a newer EOS MGM)
* Enh cs3org/reva#1944: Implement listing supported mime types in app registry
* Enh cs3org/reva#1870: Be defensive about wrongly quoted etags
* Enh cs3org/reva#1940: Reduce memory usage when uploading with S3ng storage
* Enh cs3org/reva#1888: Refactoring of the webdav code
* Enh cs3org/reva#1900: Check for illegal names while uploading or moving files
* Enh cs3org/reva#1925: Refactor listing and statting across providers for virtual views
* Fix cs3org/reva#1883: Pass directories with trailing slashes to eosclient.GenerateToken
* Fix cs3org/reva#1878: Improve the webdav error handling in the trashbin
* Fix cs3org/reva#1884: Do not send body on failed range request
* Enh cs3org/reva#1744: Add support for lightweight user types
* Fix cs3org/reva#1904: Set Content-Length to 0 when swallowing body in the datagateway
* Fix cs3org/reva#1899: Bugfix: Fix chunked uploads for new versions
* Enh cs3org/reva#1888: Refactoring of the webdav code
* Enh cs3org/reva#1887: Add "a" and "l" filter for grappa queries

https://github.com/owncloud/ocis/pull/2355
https://github.com/owncloud/ocis/pull/2295
https://github.com/owncloud/ocis/pull/2314
