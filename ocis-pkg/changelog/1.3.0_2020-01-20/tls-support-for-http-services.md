Change: Add TLS support for http services

`ocis-pkg` http services support TLS. The idea behind is setting the issuer on
phoenix's `config.json` to `https`. Or in other words, use https to access the
Kopano extension, and authenticate using an SSL certificate.

<https://github.com/owncloud/ocis-pkg/issues/19>
