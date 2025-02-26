Bugfix: We fixed the client config generation for the built in IDP

We now use the OCIS_URL to generate the web client registration configuration. It does not make sense use the OCIS_ISSUER_URL if the idp was configured to run on a different domain.

https://github.com/owncloud/ocis/pull/9770
