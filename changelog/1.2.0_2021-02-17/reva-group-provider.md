Enhancement: Enable group sharing and add config for sharing SQL driver

This PR adds config to support sharing with groups. It also introduces a
breaking change for the CS3APIs definitions since grantees can now refer to both
users as well as groups. Since we store the grantee information in a json file,
`/var/tmp/ocis/storage/shares.json`, its previous version needs to be removed as
we won't be able to unmarshal data corresponding to the previous definitions.

https://github.com/owncloud/ocis/pull/1626
https://github.com/cs3org/reva/pull/1453
