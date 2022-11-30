Enhancement: wrap metadata storage with dedicated reva gateway

We wrapped the metadata storage in a minimal reva instance with a dedicated gateway, including static storage registry, static auth registry, in memory userprovider, machine authprovider and demo permissions service. This allows us to preconfigure the service user for the ocis settings service, share and public share providers.

https://github.com/owncloud/ocis/pull/3602
https://github.com/owncloud/ocis/pull/3647
