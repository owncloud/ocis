Enhancement: Bump libre-graph-api-go

We fixed a couple of issues in libre-graph-api-go package.

* rename drive permission grantedTo to grantedToIdentities to be ms graph spec compatible.
* drive.name is a required property now.
* add group property to the identitySet.

https://github.com/owncloud/ocis/pull/5309
https://github.com/owncloud/ocis/pull/5312
