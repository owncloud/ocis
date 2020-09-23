Enhancement: Allow configuring arbitrary storage registry rules

We added a new config flag `storage-registry-rule` that can be given multiple times for the gateway to specify arbitrary storage registry rules. You can also use a comma separated list of rules in the `REVA_STORAGE_REGISTRY_RULES` environment variable.

<https://github.com/owncloud/product/issues/193>
<https://github.com/owncloud/ocis/ocis-revapull/461>
