# Changes in unreleased

## Summary

* Bugfix - Generate a random CSP-Nonce in the webapp: [#17](https://github.com/owncloud/ocis-konnectd/issues/17)
* Change - Dummy index.html is not required anymore by upstream: [#25](https://github.com/owncloud/ocis-konnectd/issues/25)
* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-konnectd/issues/1)

## Details

* Bugfix - Generate a random CSP-Nonce in the webapp: [#17](https://github.com/owncloud/ocis-konnectd/issues/17)

   https://github.com/owncloud/ocis-konnectd/issues/17
   https://github.com/owncloud/ocis-konnectd/pull/29


* Change - Dummy index.html is not required anymore by upstream: [#25](https://github.com/owncloud/ocis-konnectd/issues/25)

   The workaround was required as identifier webapp was mandatory, but we serve it from memory.
   This also introduces --disable-identifier-webapp flag.

   https://github.com/owncloud/ocis-konnectd/issues/25


* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-konnectd/issues/1)

   Just prepare an initial basic version to serve konnectd embedded into our microservice
   infrastructure in the scope of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis-konnectd/issues/1

